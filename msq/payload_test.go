package msq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalPayload(t *testing.T) {
	setup()
	defer teardown()

	assert.Equal(t, payload["example"].(string), "data")

	bytes, err := payload.Marshal()

	assert.Nil(t, err)
	assert.NotEqual(t, bytes, []byte{})
}

func TestUnMarshalPayload(t *testing.T) {
	setup()
	defer teardown()

	payload, err := payload.UnMarshal([]byte(`{"unmarshal": "testing", "numbers": [1,2,3,4]}`))

	assert.Nil(t, err)

	data, ok := (*payload)["unmarshal"]

	assert.True(t, ok)

	assert.Equal(t, data.(string), "testing")
}
