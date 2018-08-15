package msq

import (
	"context"
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

func (l *Listener) Start(handle func([]Event) bool, num int) {
	started := make(chan bool)

	if num < 1 {
		num = 1
	}

	go func() {
		if l.Running {
			panic("Cannot start the listener whilst it is already running")
		}

		defer l.cancel()

		firstTick := true

		l.interval = time.NewTicker(l.Config.Interval).C
		l.stop = make(chan bool)

		for {
			select {
			case <-l.interval:
				if !firstTick && !l.Running {
					return
				}

				if firstTick {
					l.Running = true
					started <- true
					firstTick = false
				}

				timeout := time.NewTimer(l.Config.Timeout).C

				// Go off and actually pull the events
				go func() {
					var resultValue bool
					result := make(chan bool)
					events := []Event{}

					// Depending on how many we want, that's what
					// we will pop off the queue
					for i := 0; i < num; i++ {
						event, err := l.Queue.Pop()

						if err == nil {
							events = append(events, *event)
							continue
						}
					}

					// Go off and handle those events
					go func(events []Event, handle func([]Event) bool, result chan bool) {
						result <- handle(events)
					}(events, handle, result)

					// Block on either a timeout on the handle
					// or a result from the handle.
					select {
					case <-timeout:
						for _, event := range events {
							l.Queue.ReQueue(&event)
						}
					case resultValue = <-result:
						for _, event := range events {
							if resultValue {
								l.Queue.Done(&event)
							} else {
								l.Queue.ReQueue(&event)
							}
						}

						break
					}
				}()
			case <-l.stop:
				l.Running = false
				break
			}
		}
	}()

	<-started
}

func (l *Listener) Stop() {
	l.stop <- true
}
