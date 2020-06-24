package event

import (
	"os"
	"strconv"
)

const (
	EnvKeyEventQueues  = "FLOGO_EVENT_QUEUES"
	DefaultEventQueues = 100

	EnvKeyPublishAuditEvents  = "FLOGO_PUBLISH_AUDIT_EVENTS"
	DefaultPublishAuditEvents = true
)

// PublishEventEnabled indicate the publish event enabled or not
func PublishEventEnabled() bool {
	key := os.Getenv(EnvKeyPublishAuditEvents)
	if len(key) > 0 {
		publish, _ := strconv.ParseBool(key)
		return publish
	}
	return DefaultPublishAuditEvents
}

//GetEventQueues returns the number of queues to buffer events
func GetEventQueues() int {
	numQueues := DefaultEventQueues
	queuesEnv := os.Getenv(EnvKeyEventQueues)
	if len(queuesEnv) > 0 {
		i, err := strconv.Atoi(queuesEnv)
		if err == nil {
			numQueues = i
		}
	}
	return numQueues
}
