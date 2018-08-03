package msq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanPush(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	event, err := queue.Push(payload)

	assert.Nil(t, err)
	assert.NotNil(t, event)

	assert.Equal(t, event.Namespace, queue.Config.Name)

	encodedPayload, err := payload.Marshal()

	if assert.Nil(t, err) {
		assert.Equal(t, event.Payload, string(encodedPayload))
	}
}

func TestCanPop(t *testing.T) {
	setup()
	defer teardown()
}
