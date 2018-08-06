package msq

import (
	"time"
)

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

func (q *Queue) Done(event *Event) error {
	return q.Connection.Database().Unscoped().Delete(event).Error
}

func (q *Queue) ReQueue(event *Event) error {
	return q.Connection.Database().
		Unscoped().
		Model(event).
		Update("deleted_at", nil).
		Update("retries", event.Retries+1).
		Error
}

func (q *Queue) Pop() (*Event, error) {
	event := &Event{}

	db := q.Connection.Database()

	err := db.Order("created_at desc").
		Where("retries <= ?", q.Config.MaxRetries).
		Where("namespace = ?", q.Config.Name).
		First(event).Error

	if err != nil {
		return event, err
	}

	db.Delete(event)

	return event, nil
}

func (q *Queue) Push(payload Payload) (*Event, error) {
	encodedPayload, err := payload.Marshal()

	if err != nil {
		return &Event{}, err
	}

	event := &Event{
		Namespace: q.Config.Name,
		Payload:   string(encodedPayload),
		Retries:   1,
	}

	err = q.Connection.Database().Create(event).Error

	if err != nil {
		return &Event{}, err
	}

	return event, nil
}
