package msq

var connection *Connection
var connectionConfig *ConnectionConfig
var payload Payload

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
}

func teardown() {
	connection.Close()
}
