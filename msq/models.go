package msq

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/twinj/uuid"
)

type Event struct {
	gorm.Model
	UID       string     `gorm:"type:varchar(255);primary_key"`
	Namespace string     `gorm:"type:varchar(255);index:namespace;not null"`
	Payload   string     `gorm:"type:text"`
	Retries   int        `gorm:"size:1;index:retries;default:0"`
	CreatedAt *time.Time `gorm:"index:created_at;not null"`
}

func (e *Event) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UID", uuid.NewV4())
	scope.SetColumn("CreatedAt", time.Now())

	return nil
}
