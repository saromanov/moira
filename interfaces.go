package moira

import (
	"sync"
	"time"
)

// Database implements DB functionality
type Database interface {
	FetchEvent() (*EventData, error)
	GetNotificationTrigger(id string) (TriggerData, error)
	GetTriggerTags(id string) ([]string, error)
	GetTagsSubscriptions(tags []string) ([]SubscriptionData, error)
	GetSubscription(id string) (SubscriptionData, error)
	GetContact(id string) (ContactData, error)
	GetAllContacts() ([]ContactData, error)
	SetContact(contact *ContactData) error
	AddNotification(notification *ScheduledNotification) error
	GetTriggerThrottlingTimestamps(id string) (time.Time, time.Time)
	GetTriggerEventsCount(id string, from int64) int64
	SetTriggerThrottlingTimestamp(id string, next time.Time) error
	GetNotificationsAndDelete(to int64) ([]*ScheduledNotification, error)
	GetMetricsCount() (int64, error)
	GetChecksCount() (int64, error)

	UpdateMetricsHeartbeat() error
	GetPatterns() ([]string, error)
	SaveMetrics(buffer map[string]*MatchedMetric) error
	SubscribeMetricEvents(shutdown chan bool) <-chan *MetricEvent
	GetMetricRetention(metric string) (int, error)

	GetUserSubscriptionIds(string) ([]string, error)
	GetUserContacts(string) ([]string, error)

	GetTagNames() ([]string, error)
	GetTags([]string) (map[string]TagData, error)
	GetTag(string) (TagData, error)
	SetTagMaintenance(name string, data TagData) error
	GetTagTriggerIds(tagName string) ([]string, error)
	DeleteTag(tagName string) error

	GetTriggerIds() ([]string, error)
	GetTriggerCheckIds() ([]string, int64, error)
	GetFilteredTriggerCheckIds([]string, bool) ([]string, int64, error)
	GetTrigger(string) (*Trigger, error)
	GetTriggerChecks(triggerCheckIds []string) ([]TriggerChecks, error)
	GetTriggerLastCheck(triggerId string) (*CheckData, error)
	SetTriggerLastCheck(triggerId string, checkData *CheckData) error
	SetTriggerMetricsMaintenance(triggerId string, metrics map[string]int64) error
	GetPatternTriggerIds(pattern string) ([]string, error)
	GetTriggers(triggerIds []string) ([]*Trigger, error)
	DeleteTriggerThrottling(triggerId string) error
	DeleteTrigger(triggerId string) error
	SaveTrigger(triggerId string, trigger *Trigger) error

	AddTriggerCheck(triggerId string) error

	GetEvents(string, int64, int64) ([]*EventData, error)
	PushEvent(event *EventData, ui bool) error

	DeleteContact(string, string) error
	WriteContact(contact *ContactData) error

	GetSubscriptions(subscriptionIds []string) ([]SubscriptionData, error)
	WriteSubscriptions(subscriptions []*SubscriptionData) error
	UpdateSubscription(subscription *SubscriptionData) error
	CreateSubscription(subscription *SubscriptionData) error
	DeleteSubscription(subscriptionId string, userLogin string) error

	GetNotifications(start, end int64) ([]*ScheduledNotification, int64, error)
	RemoveNotification(notificationKey string) (int64, error)

	AddPatternMetric(pattern, metric string) error
	GetPatternMetrics(pattern string) ([]string, error)
	RemovePattern(pattern string) error
	RemovePatternsMetrics(pattern []string) error
	RemovePatternWithMetrics(pattern string) error
	RemovePatternTriggers(pattern string) error

	AcquireTriggerCheckLock(triggerId string, timeout int) error
	DeleteTriggerCheckLock(triggerId string) error
	SetTriggerCheckLock(triggerId string) (*string, error)
}

// Logger implements logger abstraction
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
}

// Worker interface for implementing specified parallel workers
type Worker interface {
	Run(shutdown chan bool, wg *sync.WaitGroup)
}

// Sender interface for implementing specified contact type sender
type Sender interface {
	SendEvents(events EventsData, contact ContactData, trigger TriggerData, throttled bool) error
	Init(senderSettings map[string]string, logger Logger) error
}
