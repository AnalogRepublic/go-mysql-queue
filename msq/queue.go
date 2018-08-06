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

	now := time.Now()
	pushback := time.Now().Add(time.Millisecond * (time.Duration(event.Retries) * 100))
	retries := event.Retries + 1

	return q.Connection.Database().
		Unscoped().
		Model(event).
		Updates(map[string]interface{}{
			"deleted_at": nil,
			"created_at": pushback,
			"updated_at": now,
			"retries":    retries,
		}).
		Error
}

func (q *Queue) Pop() (*Event, error) {
	event := &Event{}

	db := q.Connection.Database()

	err := db.Order("created_at desc").
		Where("created_at <= ?", time.Now()).
		Where("retries <= ?", q.Config.MaxRetries).
		Where("namespace = ?", q.Config.Name).
		First(event).Error

	if err != nil {
		return event, err
	}

	db.Delete(event)

	return event, nil
}

func (q *Queue) Failed() ([]*Event, error) {
	events := []*Event{}

	db := q.Connection.Database()

	err := db.Unscoped().Order("created_at desc").
		Where("namespace = ?", q.Config.Name).
		Find(&events).
		Error

	if err != nil {
		return events, err
	}

	return events, nil
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
