package msq

import (
	"time"
)

var connection *Connection
var connectionConfig *ConnectionConfig
var payload Payload
var queueConfig *QueueConfig
var listenerConfig ListenerConfig

func setup() {
	connectionConfig = &ConnectionConfig{
		Type:     "sqlite",
		Database: "../test.db",
		Charset:  "utf8",
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
		MaxRetries: 4,
		MessageTTL: 5 * time.Minute,
	}

	listenerConfig = ListenerConfig{
		Interval: 100 * time.Millisecond,
		Timeout:  250 * time.Millisecond,
	}
}

func teardown() {
	connection.Database().DropTable(&Event{})

	connection.Close()
}
