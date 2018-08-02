package msq

type ConnectionConfig struct {
	Type     string
	Host     string
	Username string
	Password string
	Database string
}

type Connection struct {
}

func (c *Connection) Attempt() error {

}

func Connect(config ConnectionConfig) (*Queue, error) {
	connection := &Connection{
		Config: config,
	}

	err := connection.Attempt()

	if err != nil {
		return &Queue{}, err
	}

	queue := &Queue{
		Connection: connection,
	}

	return queue, nil
}
