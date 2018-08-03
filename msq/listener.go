package msq

import "time"

type ListenerConfig struct {
	Interval time.Duration
	Timeout  time.Duration
}

type Listener struct {
	Started bool
	Queue   Queue
	Config  ListenerConfig
}

func (l *Listener) Start(handle func(Event) bool) error {
	l.Started = true
	return nil
}

func (l *Listener) Stop() error {
	l.Started = false
	return nil
}
