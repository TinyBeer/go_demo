package common

import (
	"xorm.io/xorm"
	"xorm.io/xorm/log"

	// 没有引入会报错：sql: unknown driver "mysql" (forgotten import?)
	_ "github.com/go-sql-driver/mysql"
)

func Engine() *xorm.Engine {
	engine, _ := xorm.NewEngine("mysql", "root:123456@/test?charset=utf8mb4")
	engine.ShowSQL(true)
	engine.Logger().SetLevel(log.LOG_DEBUG)
	return engine
}
