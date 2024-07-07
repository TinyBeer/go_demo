package common

import (
	"xorm.io/xorm"

	// 没有引入会报错：sql: unknown driver "mysql" (forgotten import?)
	_ "github.com/go-sql-driver/mysql"
)

func NewEngineGroup() *xorm.EngineGroup {
	conns := []string{
		"root:123456@/test?charset=utf8mb4",  // 第一个默认是master
		"root:123456@/test1?charset=utf8mb4", // 第二个开始都是slave
		"root:123456@/test2?charset=utf8mb4",
	}

	eg, err := xorm.NewEngineGroup("mysql", conns, xorm.RandomPolicy())
	if err != nil {
		panic(err)
	}
	return eg
}

func NewEngineGroup2() *xorm.EngineGroup {
	var err error
	master, err := xorm.NewEngine("mysql", "root:123456@/test?charset=utf8mb4")
	if err != nil {
		panic(err)
	}

	slave1, err := xorm.NewEngine("mysql", "root:123456@/test1?charset=utf8mb4")
	if err != nil {
		panic(err)
	}

	slave2, err := xorm.NewEngine("mysql", "root:123456@/test2?charset=utf8mb4")
	if err != nil {
		panic(err)
	}

	slaves := []*xorm.Engine{slave1, slave2}
	eg, err := xorm.NewEngineGroup(master, slaves)
	if err != nil {
		panic(err)
	}
	eg.SetPolicy(xorm.LeastConnPolicy())
	return eg
}
