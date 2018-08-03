package msq

import "time"

type QueueConfig struct {
	Name       string
	MaxRetries int64
	MessageTTL time.Duration
}

type Queue struct {
	Connection *Connection
	Config     *QueueConfig
}

func (q *Queue) Configure(config *QueueConfig) {
	q.Config = config
}

func (q *Queue) Done(event *Event) {

}

func (q *Queue) ReQueue(event *Event) {

}

func (q *Queue) Listen(handle func(Event) bool, config ListenerConfig) (*Listener, error) {
	listener := &Listener{
		Queue:  *q,
		Config: config,
	}

	err := listener.Start(handle)

	return listener, err
}

func (q *Queue) Pop() (*Event, error) {
	return &Event{}, nil
}

func (q *Queue) Push(payload Payload) (*Event, error) {
	encodedPayload, err := payload.Marshal()

	if err != nil {
		return &Event{}, err
	}

	event := &Event{
		Namespace: q.Config.Name,
		Payload:   string(encodedPayload),
	}

	q.Connection.Database().Create(event)

	return event, nil
}
