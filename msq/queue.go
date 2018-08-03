package msq

type QueueConfig struct {
	Name       string
	MaxRetries int64
	MessageTTL int64
}

type Queue struct {
	Connection *Connection
	Config     *QueueConfig
}

func (q *Queue) Configure(config *QueueConfig) {
	q.Config = config
}

func (q *Queue) Done(message *Message) {

}

func (q *Queue) ReQueue(message *Message) {

}

func (q *Queue) Listen(handle func(Message) bool, config ListenerConfig) (*Listener, error) {
	listener := &Listener{
		Queue:  *q,
		Config: config,
	}

	err := listener.Start(handle)

	return listener, err
}

func (q *Queue) Pop() (*Message, error) {
	return &Message{}, nil
}

func (q *Queue) Push(payload Payload) (*Message, error) {
	return &Message{}, nil
}
