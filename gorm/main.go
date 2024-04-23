package main

import (
	"fmt"
	"learn_gorm/dao"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func init() {
	username := "root"
	password := "123456"
	host := "192.168.56.101"
	port := 3306
	dbname := "gorm"
	timeout := 10

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%ds",
		username, password, host, port, dbname, timeout)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		// SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:         "gorm_",
			SingularTable:       true,
			NameReplacer:        nil,
			NoLowerCase:         false,
			IdentifierMaxLength: 0,
		},
	})
	if err != nil {
		panic(err)
	}

	DB = db.Debug()
}

func main() {
	dao := dao.NewUserDao(DB)

	err := dao.CreateUser("tom", "hello")
	if err != nil {
		fmt.Println(err)
		return
	}

	u, err := dao.GetUserByName("tom")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(u)

	time.Sleep(time.Second * 3)

	err = dao.UpdatePassword(u.Id, "newpwd")
	if err != nil {
		fmt.Println(err)
		return
	}

	u, err = dao.GetUserById(u.Id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(u)

	time.Sleep(time.Second * 5)
	err = dao.DeleteUserById(u.Id)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = dao.GetUserById(u.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("用户不存在")
		} else {
			fmt.Println(err)
		}
	}
}
