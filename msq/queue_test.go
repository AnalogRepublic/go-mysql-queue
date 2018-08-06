package msq

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	event, err := queue.Push(payload)

	if assert.Nil(t, err) {
		assert.NotNil(t, event)

		assert.Equal(t, event.Namespace, queue.Config.Name)

		encodedPayload, err := payload.Marshal()

		if assert.Nil(t, err) {
			assert.Equal(t, event.Payload, string(encodedPayload))
		}
	}
}

func TestPop(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	_, err = queue.Push(payload)

	if assert.Nil(t, err) {
		event, err := queue.Pop()

		if assert.Nil(t, err) {
			if assert.NotEqual(t, event.UID, "", "UID should not be empty") {
				encodedPayload, err := payload.Marshal()

				if assert.Nil(t, err) {
					assert.Equal(t, event.Payload, string(encodedPayload), "Payload should match")
				}
			}
		}
	}
}

func TestDone(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	event, err := queue.Pop()

	if assert.Nil(t, err, "There should be an event in the queue") {
		err := queue.Done(event)

		assert.Nil(t, err, "We should be able to remove the record")

		err = connection.Database().Where("uid = ?", event.UID).First(&Event{}).Error

		assert.NotNil(t, err, "We want the record to be missing as it should be removed")
	}
}

func TestReQueue(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	originalEvent, err := queue.Push(payload)

	if assert.Nil(t, err) {
		event, err := queue.Pop()

		if assert.Nil(t, err, "There should be an event in the queue") {
			// Simulate waiting
			time.Sleep(listenerConfig.Interval)

			err := queue.ReQueue(event)
			assert.Nil(t, err, "We should have no problem re-queuing the event")

			// Simulate waiting
			time.Sleep(listenerConfig.Interval)

			newEvent, err := queue.Pop()

			if assert.Nil(t, err, "We should find a requeued event in the queue") {
				assert.Equal(t, event.UID, newEvent.UID, "We should get back the same event")
				assert.NotEqual(t, originalEvent.Retries, newEvent.Retries, "The retries should have increased")
				assert.NotEqual(t, originalEvent.UpdatedAt.UnixNano(), newEvent.UpdatedAt.UnixNano(), "The updated at timestamps should be different")
				assert.NotEqual(t, originalEvent.CreatedAt.UnixNano(), newEvent.CreatedAt.UnixNano(), "The created at timestamps should be different")
			}
		}
	}
}

func TestFailed(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	originalEvent, err := queue.Push(payload)

	if assert.Nil(t, err) {
		event, err := queue.Pop()

		if assert.Nil(t, err, "There should be an event in the queue") {
			// Simulate waiting
			time.Sleep(listenerConfig.Interval)

			event.Retries = 3

			err := queue.ReQueue(event)
			assert.Nil(t, err, "We should have no problem re-queuing the event")

			// Simulate waiting
			time.Sleep(listenerConfig.Interval)

			failedEvents, err := queue.Failed()

			if assert.Nil(t, err, "We should get a list of a failed events") {
				assert.Equal(t, originalEvent.UID, failedEvents[0].UID)
				queue.Done(failedEvents[0])
			}
		}
	}
}
