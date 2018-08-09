package msq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventPayload(t *testing.T) {
	event := Event{
		Payload: `{"example":"json", "payload":"4tests"}`,
	}

	assert.NotEmpty(t, event.Payload)

	payload, err := event.GetPayload()

	if assert.Nil(t, err) {
		assert.Equal(t, payload["example"], "json")
		assert.Equal(t, payload["payload"], "4tests")
	}
}
