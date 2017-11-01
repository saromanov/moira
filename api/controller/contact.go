package controller

import (
	"time"

	"github.com/go-graphite/carbonapi/date"
	"github.com/satori/go.uuid"

	"github.com/moira-alert/moira"
	"github.com/moira-alert/moira/api"
	"github.com/moira-alert/moira/api/dto"
)

// GetAllContacts gets all moira contacts
func GetAllContacts(database moira.Database) (*dto.ContactList, *api.ErrorResponse) {
	contacts, err := database.GetAllContacts()
	if err != nil {
		return nil, api.ErrorInternalServer(err)
	}
	contactsList := dto.ContactList{
		List: contacts,
	}
	return &contactsList, nil
}

// CreateContact creates new notification contact for current user
func CreateContact(database moira.Database, contact *dto.Contact, userLogin string) *api.ErrorResponse {
	id := uuid.NewV4().String()
	contactData := &moira.ContactData{
		ID:    id,
		User:  userLogin,
		Type:  contact.Type,
		Value: contact.Value,
	}
	if err := database.SaveContact(contactData); err != nil {
		return api.ErrorInternalServer(err)
	}
	contact.User = &userLogin
	contact.ID = &id
	return nil
}

// RemoveContact deletes notification contact for current user and remove contactID from all subscriptions
func RemoveContact(database moira.Database, contactID string, userLogin string) *api.ErrorResponse {
	subscriptionIDs, err := database.GetUserSubscriptionIDs(userLogin)
	if err != nil {
		return api.ErrorInternalServer(err)
	}

	subscriptions, err := database.GetSubscriptions(subscriptionIDs)
	if err != nil {
		return api.ErrorInternalServer(err)
	}

	subscriptionsWithDeletingContact := make([]*moira.SubscriptionData, 0)

	for _, subscription := range subscriptions {
		if subscription == nil {
			continue
		}
		for i, contact := range subscription.Contacts {
			if contact == contactID {
				subscription.Contacts = append(subscription.Contacts[:i], subscription.Contacts[i+1:]...)
				subscriptionsWithDeletingContact = append(subscriptionsWithDeletingContact, subscription)
				break
			}
		}
	}

	if err := database.RemoveContact(contactID); err != nil {
		return api.ErrorInternalServer(err)
	}

	if err := database.SaveSubscriptions(subscriptionsWithDeletingContact); err != nil {
		return api.ErrorInternalServer(err)
	}

	return nil
}

// TestContact push test notification to verify the correct contact settings
func TestContact(database moira.Database, contactID string) *api.ErrorResponse {
	var value float64 = 1
	eventData := &moira.NotificationEvent{
		ContactID: contactID,
		Metric:    "Test.metric.value",
		Value:     &value,
		OldState:  "TEST",
		State:     "TEST",
		Timestamp: int64(date.DateParamToEpoch("now", "", time.Now().Add(-24*time.Hour).Unix(), time.UTC)),
	}
	if err := database.PushNotificationEvent(eventData, false); err != nil {
		return api.ErrorInternalServer(err)
	}
	return nil
}
