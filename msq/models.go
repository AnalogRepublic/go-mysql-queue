package msq

import "github.com/jinzhu/gorm"

type Event struct {
	gorm.Model
	UID       string `gorm:"type:varchar(255);primary_key"`
	Namespace string `gorm:"type:varchar(255);index:namespace;not null"`
}
