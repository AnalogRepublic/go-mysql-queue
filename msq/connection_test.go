package msq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionConnect(t *testing.T) {
	setup()
	defer teardown()

	connection = &Connection{
		Config: *connectionConfig,
	}

	err := connection.Attempt()

	assert.Nil(t, err)
}

func TestConnectionMigrateDatabase(t *testing.T) {
	setup()
	defer teardown()

	connection = &Connection{
		Config: *connectionConfig,
	}

	err := connection.Attempt()

	assert.Nil(t, err)

	err = connection.SetupDatabase()

	assert.Nil(t, err)
}

func TestConnectionClose(t *testing.T) {
	setup()

	connection = &Connection{
		Config: *connectionConfig,
	}

	err := connection.Attempt()

	assert.Nil(t, err)

	err = connection.Close()

	assert.Nil(t, err)
}

func TestConnect(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	assert.Nil(t, err)

	if assert.NotNil(t, queue) {
		assert.Equal(t, queue.Connection.Config.Type, config.Type)
		assert.Equal(t, queue.Connection.Config.Database, config.Database)
	}
}
