package dao

import (
	"learn_gorm/domain"
	"learn_gorm/model"

	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
	return &UserDao{
		db: db,
	}
}

func (u *UserDao) CreateUser(name string, password string) error {
	return u.db.Create(&model.User{
		Name:     name,
		Password: password,
	}).Error
}

func (u *UserDao) GetUserById(id uint) (*model.User, error) {
	user := &model.User{}
	err := u.db.Take(user, id).Error
	return user, err
}

func (u *UserDao) GetUserByName(name string) (*model.User, error) {
	user := &model.User{}
	err := u.db.Take(user, "name = ?", name).Error
	return user, err
}

func (u *UserDao) UpdatePassword(id uint, password string) error {
	return u.db.Model(&model.User{}).Where("id = ?", id).Update("password", password).Error
}

func (u *UserDao) DeleteUserById(id uint) error {
	return u.db.Delete(&model.User{}, id).Error
}

func (u *UserDao) DeleteUserByName(name string) error {
	return u.db.Delete(&model.User{}, "name = ?", name).Error
}

var _ domain.IUserDao = &UserDao{}
