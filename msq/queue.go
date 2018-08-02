package msq

type Queue struct {
	Connection *Connection
	Config     *QueueConfig
}

func (q *Queue) Configure(config *QueueConfig) {
	q.Config = config
}

func (q *Queue) SetupDatabase() error {

}

func (q *Queue) Listen(handle func(Message) bool, config ListenerConfig) (*Listener, error) {
	listener := &Listener{
		Queue:  q,
		Config: config,
	}

	err := listener.Start(handle)

	return listener, err
}

func (q *Queue) Pop() (*Message, error) {

}

func (q *Queue) Push(payload Payload) (*Message, error) {

}
