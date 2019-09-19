package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/go-redis/redis"
	"github.com/jackyczj/NoGhost/Auth"
	cacheClient "github.com/jackyczj/NoGhost/cache"
	"github.com/jackyczj/NoGhost/store"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:   Auth.Skipper,
		Validator: Auth.Validator,
	}))
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Mongo Break ,maybe you should check it out.")
			}
		}()
		err := store.InitDB().Client().Ping(ctx, readpref.Primary())
		if err != nil {
			panic(err)
		}
		fmt.Println("MongoDB OK!")
	}()
	defer store.Client.Close()
	e.Use(middleware.CORS())
	op := redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Redis down , please test your env.")
			}
		}()
		cacheClient.InitCache()
		cacheClient.Rdb = redis.NewClient(&op)
		_, err := cacheClient.Rdb.Ping().Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("Redis OK!")
	}()

	e.Logger.Fatal("Service start at port:", e.Start(":2333"))
}
