package main

import (
	"log"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {

	v := viper.New()
	// 没有对账号密码的支持
	v.AddRemoteProvider("etcd3", "http://192.168.56.101:2379", "/config.json")
	// v.AddSecureRemoteProvider()
	v.SetConfigType("json")
	err := v.ReadRemoteConfig()
	if err != nil {
		log.Fatalln(err)
	}
	v.SetConfigFile("./config.json")

	v.WriteConfig()

}
