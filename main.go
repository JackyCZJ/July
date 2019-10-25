package main

import (
	"context"
	"os"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/spf13/pflag"

	"github.com/jackyczj/July/config"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	cacheClient "github.com/jackyczj/July/cache"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)
}

var (
	cfg = pflag.StringP("config", "c", "", "July config file path.")
)

func main() {
	pflag.Parse()
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}
	e := echo.New()
	Load(e)
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
		_, err := cacheClient.Cluster.Ping().Result()
		if err != nil {
			panic(err)
		}
	}()
	go func(e *echo.Echo) {
		e.AutoTLSManager.Cache = autocert.DirCache("/conf")
		e.Logger.Fatal("TLS service start at port:", e.StartAutoTLS(":443"))
	}(e)
	e.Logger.Fatal("Service start at port:", e.Start(":2333"))

}
