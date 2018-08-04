package msq

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartStop(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	queuedEvent, err := queue.Push(payload)

	if assert.Nil(t, err) {
		listener := &Listener{
			Queue:  *queue,
			Config: listenerConfig,
		}

		go listener.Start(func(event Event) bool {
			assert.Equal(t, queuedEvent.UID, event.UID)
			return true
		})

		time.Sleep(time.Second)

		assert.Equal(t, listener.Config.Interval, listenerConfig.Interval)
		assert.Equal(t, listener.Config.Timeout, listenerConfig.Timeout)

		assert.True(t, listener.Started, "The listener should be started")

		err := listener.Stop()
		assert.Nil(t, err, "We should not error when trying to stop")

		time.Sleep(time.Second)
		assert.False(t, listener.Started, "The listener should now be stopped")
	}
}

func TestHandleFail(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	queuedEvent, err := queue.Push(payload)

	if assert.Nil(t, err) {
		listener := &Listener{
			Queue:  *queue,
			Config: listenerConfig,
		}

		go listener.Start(func(event Event) bool {
			assert.Equal(t, queuedEvent.UID, event.UID)
			return false
		})

		time.Sleep(time.Second)

		assert.Equal(t, listener.Config.Interval, listenerConfig.Interval)
		assert.Equal(t, listener.Config.Timeout, listenerConfig.Timeout)

		assert.True(t, listener.Started, "The listener should be started")

		poppedEvent, err := queue.Pop()

		if assert.Nil(t, err, "We should get an event back as it should've been re-queued") {
			assert.Equal(t, poppedEvent.UID, queuedEvent.UID)
			queue.Done(poppedEvent)
		}

		queue.Done(poppedEvent)

		err = listener.Stop()
		assert.Nil(t, err, "We should not error when trying to stop")

		time.Sleep(time.Second)
		assert.False(t, listener.Started, "The listener should now be stopped")
	}
}

func TestHandleTimeout(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	queuedEvent, err := queue.Push(payload)

	if assert.Nil(t, err) {
		listener := &Listener{
			Queue:  *queue,
			Config: listenerConfig,
		}

		go listener.Start(func(event Event) bool {
			time.Sleep(time.Second)
			assert.Equal(t, queuedEvent.UID, event.UID)
			return true
		})

		time.Sleep(time.Second)

		assert.Equal(t, listener.Config.Interval, listenerConfig.Interval)
		assert.Equal(t, listener.Config.Timeout, listenerConfig.Timeout)

		assert.True(t, listener.Started, "The listener should be started")

		poppedEvent, err := queue.Pop()

		if assert.Nil(t, err, "We should get an event back as it should've been re-queued") {
			assert.Equal(t, poppedEvent.UID, queuedEvent.UID)
			queue.Done(poppedEvent)
		}

		err = listener.Stop()
		assert.Nil(t, err, "We should not error when trying to stop")

		time.Sleep(time.Second)
		assert.False(t, listener.Started, "The listener should now be stopped")
	}
}
