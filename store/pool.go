package store

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var rdb *redis.Client

func init(){
	opt := &redis.Options{
		Addr:         ":6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	}
	rdb = redis.NewClient(opt)
	if rdb.Ping().Err() != nil{
		fmt.Println(rdb.Ping().Err().Error())
	}
}