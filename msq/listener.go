package msq

import (
	"time"
)

type ListenerConfig struct {
	Interval time.Duration
	Timeout  time.Duration
}

type Listener struct {
	Started         bool
	Queue           Queue
	Config          ListenerConfig
	interval        <-chan time.Time
	listenerRunning chan bool
}

func (l *Listener) Start(handle func(Event) bool) {
	if l.Started {
		panic("Cannot start the listener whilst it is already running")
	}

	l.Started = true

	l.interval = time.NewTicker(l.Config.Interval).C
	l.listenerRunning = make(chan bool)

	for {
		select {
		case <-l.interval:
			if !l.Started {
				return
			}

			go func() {
				event, err := l.Queue.Pop()

				if err == nil {
					timeout := time.NewTimer(l.Config.Timeout).C

					var resultValue bool
					result := make(chan bool)

					go func(event Event, handle func(Event) bool, result chan bool) {
						result <- handle(event)
					}(*event, handle, result)

					select {
					case <-timeout:
						l.Queue.ReQueue(event)
					case resultValue = <-result:
						if resultValue {
							l.Queue.Done(event)
						} else {
							l.Queue.ReQueue(event)
						}
					}
				}
			}()
		case <-l.listenerRunning:
			l.Started = false
			return
		}
	}
}

func (l *Listener) Stop() error {
	l.listenerRunning <- true
	return nil
}
