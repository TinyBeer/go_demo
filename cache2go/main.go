package main

import (
	"fmt"
	"time"

	"github.com/muesli/cache2go"
)

// Keys & values in cache2go can be of arbitrary types, e.g. a struct.
type myStruct struct {
	text     string
	moreData []byte
}

func main() {
	cache := cache2go.Cache("my_cache")
	item := cache.Add("hello", time.Second, "world")
	item.SetAboutToExpireCallback(func(key interface{}) {
		fmt.Println(key, "deleted")
	})
	time.Sleep(time.Second * 2)

	item = cache.Add("hello", time.Second, "world")
	item.AddAboutToExpireCallback(func(key interface{}) {
		fmt.Println(key, "deleted")
	})
	item.AddAboutToExpireCallback(func(key interface{}) {
		fmt.Println(key, "deleted2")
	})
	cache.Delete("hello")
}
