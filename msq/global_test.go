package msq

import (
	"time"
)

var connection *Connection
var connectionConfig *ConnectionConfig
var payload Payload
var queueConfig *QueueConfig

func setup() {
	connectionConfig = &ConnectionConfig{
		Type:     "sqlite",
		Database: "../test.db",
	}

	payload = Payload{
		"example": "data",
		"is": map[string]string{
			"being": "shown",
		},
		"here": []int{1, 2, 3, 4},
	}

	queueConfig = &QueueConfig{
		Name:       "testing",
		MaxRetries: 3,
		MessageTTL: 5 * time.Minute,
	}
}

func teardown() {
	connection.Database().DropTable(&Event{})

	connection.Close()
}
