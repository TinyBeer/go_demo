package main

import (
	"github.com/spf13/viper"
)

func main() {
	username := "root"
	password := "123456"
	host := "192.168.56.101"
	port := 3306
	dbname := "gorm"
	timeout := 10

	v := viper.New()
	v.SetConfigFile("./config.yaml")
	v.Set("mysql.username", username)
	v.Set("mysql.password", password)
	v.Set("mysql.host", host)
	v.Set("mysql.port", port)
	v.Set("mysql.dbname", dbname)
	v.Set("mysql.timeout", timeout)
	v.WriteConfig()
}
