package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	Pool *redis.Pool
)

func init() {
	redisHost := ":6379"
	Pool = newPool(redisHost)
	close()
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func close() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}

func Set(key, value string) error {
	conn := Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("SET", key, value); err != nil {
		return err
	}

	return nil
}

func Get(key string) error {
	conn := Pool.Get()
	defer conn.Close()

	reply, err := conn.Do("GET", key)
	if err != nil {
		return err
	}

	value, err := redis.Bytes(reply, nil)
	if err != nil {
		return err
	}

	fmt.Println(string(value))
	return nil
}

func main() {
	key, value := "name", "hxia"
	if err := Set(key, value); err != nil {
		panic(err)
	}

	if err := Get(key); err != nil {
		panic(err)
	}
}
