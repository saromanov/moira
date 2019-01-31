package webhook

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/moira-alert/moira"
)

// Sender implements moira sender interface via webhook
type Sender struct {
	url          string
	user         string
	password     string
	headers      map[string]string
	client       *http.Client
}

// Init read yaml config
func (sender *Sender) Init(senderSettings map[string]string, logger moira.Logger, location *time.Location, dateTimeFormat string) error {

	if senderSettings["name"] == "" {
		return fmt.Errorf("required name for sender type webhook")
	}

	sender.url = senderSettings["url"]
	if sender.url == "" {
		return fmt.Errorf("can not read url from config")
	}

	sender.user, sender.password = senderSettings["user"], senderSettings["password"]

	senderHeaders := make(map[string]string)

	if headers, ok := senderSettings["headers"]; ok {
		err := yaml.Unmarshal([]byte(headers), senderHeaders)
		if err != nil {
			return fmt.Errorf("can not read headers from config: %s", err.Error())
		}
		sender.headers = senderHeaders
	}

	timeout := 30

	if timeoutRaw, ok := senderSettings["timeout"]; ok {
		var err error
		timeout, err = strconv.Atoi(timeoutRaw)
		if err != nil {
			return fmt.Errorf("can not read timeout from config: %s", err.Error())
		}
	}

	tr := &http.Transport{DisableKeepAlives: true}
	sender.client = &http.Client{Timeout: time.Duration(timeout) * time.Second, Transport: tr}
	return nil
}

// SendEvents implements Sender interface Send
func (sender *Sender) SendEvents(events moira.NotificationEvents, contact moira.ContactData, trigger moira.TriggerData, plot []byte, throttled bool) error {
	request, err := sender.buildRequest(events, contact, trigger, plot, throttled)
	if request != nil {
		defer request.Body.Close()
	}

	if err != nil {
		return fmt.Errorf("failed to build request body: %s", err.Error())
	}

	if sender.user != "" && sender.password != "" {
		request.SetBasicAuth(sender.user, sender.password)
	}

	for k, v := range sender.headers {
		request.Header.Set(k, v)
	}

	response, err := sender.client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return fmt.Errorf("failed to perform request: %s", err.Error())
	}

	if !isAllowedResponseCode(response.StatusCode) {
		var serverResponse string
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			serverResponse = fmt.Sprintf("failed to read response body: %s", err.Error())
		} else {
			serverResponse = string(responseBody)
		}
		return fmt.Errorf("invalid status code: %d, server response: %s", response.StatusCode, serverResponse)
	}

	return nil
}

func isAllowedResponseCode(responseCode int) bool {
	return (responseCode >= 200) && (responseCode <= 299)
}
