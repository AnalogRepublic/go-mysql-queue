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

		assert.Equal(t, listener.Config.Interval, listenerConfig.Interval)
		assert.Equal(t, listener.Config.Timeout, listenerConfig.Timeout)

		ctx := listener.Context()

		listener.Start(func(event Event) bool {
			assert.Equal(t, queuedEvent.UID, event.UID)
			return true
		})

		go func() {
			assert.True(t, listener.Running, "The listener should be running")

			time.Sleep(time.Second)
			listener.Stop()
		}()

		select {
		case <-ctx.Done():
			assert.False(t, listener.Running, "The listener should no longer be running")
		}
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

		assert.Equal(t, listener.Config.Interval, listenerConfig.Interval)
		assert.Equal(t, listener.Config.Timeout, listenerConfig.Timeout)

		ctx := listener.Context()

		listener.Start(func(event Event) bool {
			assert.Equal(t, queuedEvent.UID, event.UID)
			return true
		})

		go func() {
			assert.True(t, listener.Running, "The listener should be started")
			time.Sleep(time.Second)

			poppedEvent, err := queue.Pop()

			if assert.Nil(t, err, "We should get an event back as it should've been re-queued") {
				assert.Equal(t, poppedEvent.UID, queuedEvent.UID)
				queue.Done(poppedEvent)
			}

			queue.Done(poppedEvent)

			listener.Stop()
		}()

		select {
		case <-ctx.Done():
			assert.False(t, listener.Running, "The listener should no longer be running")
		}
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

		assert.Equal(t, listener.Config.Interval, listenerConfig.Interval)
		assert.Equal(t, listener.Config.Timeout, listenerConfig.Timeout)

		ctx := listener.Context()

		listener.Start(func(event Event) bool {
			time.Sleep(time.Second)
			assert.Equal(t, queuedEvent.UID, event.UID)
			return true
		})

		go func() {
			assert.True(t, listener.Running, "The listener should be started")
			time.Sleep(time.Second)

			poppedEvent, err := queue.Pop()

			if assert.Nil(t, err, "We should get an event back as it should've been re-queued") {
				assert.Equal(t, poppedEvent.UID, queuedEvent.UID)
				queue.Done(poppedEvent)
			}

			listener.Stop()
		}()

		select {
		case <-ctx.Done():
			assert.False(t, listener.Running, "The listener should no longer be running")
		}
	}
}
