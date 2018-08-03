package msq

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	UID       string     `gorm:"type:varchar(255);primary_key"`
	Namespace string     `gorm:"type:varchar(255);index:namespace;not null"`
	Name      string     `gorm:"type:varchar(255);index:name"`
	Payload   string     `gorm:"type:text"`
	Retries   int        `gorm:"size:1;index:retries"`
	CreatedAt *time.Time `gorm:"index:created_at;not null"`
}
