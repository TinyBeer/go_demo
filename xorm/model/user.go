package model

import "time"

type User struct {
	ID      int64
	Name    string
	Salt    string
	Age     int
	Level   int
	Passwd  string    `xorm:"varchar(200)"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func (*User) TableName() string {
	return "my_user"
}

func (*User) Charset() string {
	return "utf8mb4"
}

func (*User) StoreEngine() string {
	return "innodb"
}
