package model

import (
	"gorm.io/plugin/soft_delete"
)

type User struct {
	Id        uint                  `gorm:"size:10"`
	Name      string                `gorm:"size:10;not null"`
	Password  string                `gorm:"size:20;not null"`
	CreatedAt int64                 `gorm:"autoCreateTime"`
	UpdatedAt int64                 `gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:udx_name"`
}
