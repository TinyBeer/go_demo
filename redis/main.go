package main

import (
	"fmt"

	redis "github.com/gomodule/redigo/redis"
)

func TestRedisPool() {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				"192.168.10.103:6379",
				redis.DialPassword("123456"),
			)
		},
		MaxIdle:     8,
		MaxActive:   0,
		IdleTimeout: 100,
	}
	conn := pool.Get()
	defer conn.Close()

	fmt.Println(conn.Do("PING"))
}

func TestRedisString() {
	c, err := redis.Dial("tcp", "192.168.10.103:6379", redis.DialPassword("123456"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	_, err = c.Do("Set", "key", 999)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := redis.String(c.Do("Get", "key"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

func TestRedisHash() {
	c, err := redis.Dial("tcp", "192.168.10.103:6379", redis.DialPassword("123456"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	_, err = c.Do("HSet", "user01", "name", "john")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = c.Do("HSet", "user01", "age", 18)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := redis.String(c.Do("HGet", "user01", "name"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("res=", res)

	res2, err := redis.String(c.Do("HGet", "user01", "age"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("res2=", res2)
}

func TestRedisHash2() {
	c, err := redis.Dial("tcp", "192.168.10.103:6379", redis.DialPassword("123456"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	_, err = c.Do("HMSet", "user01", "name", "tom", "age", 88)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := redis.Strings(c.Do("HMGet", "user01", "name", "age"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("res=", res)

}

func main() {
	TestRedisPool()
	// TestRedisString()
	// TestRedisHash()
	// TestRedisHash2()
}
