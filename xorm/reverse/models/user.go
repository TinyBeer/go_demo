package models

import (
	"time"
)

type User struct {
	ID      int64     `xorm:"not null pk autoincr BIGINT(20)"`
	Name    string    `xorm:"VARCHAR(255)"`
	Salt    string    `xorm:"VARCHAR(255)"`
	Age     int       `xorm:"INT(11)"`
	Level   int       `xorm:"INT(11)"`
	Passwd  string    `xorm:"VARCHAR(200)"`
	Created time.Time `xorm:"DATETIME"`
	Updated time.Time `xorm:"DATETIME"`
}

func (m *User) TableName() string {
	return "user"
}
