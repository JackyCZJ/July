package main

import (
	"context"
	"os"
	"time"

	"github.com/spf13/pflag"

	"github.com/jackyczj/July/config"

	"github.com/jackyczj/July/handler/user"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	Auth "github.com/jackyczj/July/auth"
	cacheClient "github.com/jackyczj/July/cache"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)
}

var (
	cfg = pflag.StringP("config", "c", "", "ovn-client config file path.")
)

func main() {
	pflag.Parse()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:   Auth.Skipper,
		Validator: Auth.Validator,
	}))
	// init config

	if err := config.Init(*cfg); err != nil {
		panic(err)
	}
	e.POST("/login", user.Login)
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Mongo Break ,maybe you should check it out.")
			}
		}()
		err := store.InitDB().Client().Ping(ctx, readpref.Primary())
		if err != nil {
			panic(err)
		}
		log.Println("MongoDB OK!")
	}()
	defer store.Client.Close()
	e.Use(middleware.CORS())
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Redis down , please test your env.")
			}
		}()
		cacheClient.InitCache()
		_, err := cacheClient.Rdb.Ping().Result()
		if err != nil {
			panic(err)
		}
		log.Println("Redis OK!")
	}()

	e.Logger.Fatal("Service start at port:", e.Start(":2333"))
}
