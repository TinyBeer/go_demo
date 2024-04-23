package domain

import "learn_gorm/model"

type IUserDao interface {
	CreateUser(name string, password string) error
	GetUserById(id uint) (*model.User, error)
	GetUserByName(name string) (*model.User, error)
	UpdatePassword(id uint, password string) error
	DeleteUserById(id uint) error
	DeleteUserByName(name string) error
}
