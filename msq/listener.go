package msq

type ListenerConfig struct {
	Interval int64
	Timeout  int64
}

type Listener struct {
	Queue  Queue
	Config ListenerConfig
}

func (l *Listener) Start(handle func(Message) bool) error {

}
