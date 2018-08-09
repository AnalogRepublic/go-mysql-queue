package msq

import (
	"github.com/jinzhu/gorm"
	"github.com/twinj/uuid"
)

type Event struct {
	gorm.Model
	UID       string `gorm:"type:varchar(255);index:uid"`
	Namespace string `gorm:"type:varchar(255);index:namespace;not null"`
	Payload   string `gorm:"type:text"`
	Retries   int    `gorm:"size:1;index:retries;default:0"`
}

func (e *Event) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UID", uuid.NewV4().String())
	return nil
}

func (e *Event) GetPayload() (Payload, error) {
	payload, err := payload.UnMarshal([]byte(e.Payload))

	if err != nil {
		return Payload{}, err
	}

	return *payload, nil
}
