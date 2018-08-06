package msq

import (
	"context"
	"sync"
	"time"
)

type ListenerConfig struct {
	Interval time.Duration
	Timeout  time.Duration
}

type Listener struct {
	Running  bool
	Queue    Queue
	Config   ListenerConfig
	interval <-chan time.Time
	stop     chan bool
	ctx      context.Context
	cancel   func()
}

func (l *Listener) Context() context.Context {
	if l.ctx == nil {
		l.ctx, l.cancel = context.WithCancel(context.Background())
	}

	return l.ctx
}

func (l *Listener) Start(handle func(Event) bool) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		if l.Running {
			panic("Cannot start the listener whilst it is already running")
		}

		defer l.cancel()

		l.Running = true
		l.interval = time.NewTicker(l.Config.Interval).C
		l.stop = make(chan bool)

		wg.Done()

		for {
			select {
			case <-l.interval:
				if !l.Running {
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
			case <-l.stop:
				l.Running = false
				break
			}
		}
	}()

	wg.Wait()
}

func (l *Listener) Stop() {
	l.stop <- true
}
