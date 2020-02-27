package main

import (
	"context"
	"time"

	cacheClient "github.com/jackyczj/July/cache"
	"github.com/jackyczj/July/config"
	"github.com/jackyczj/July/log"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	cfg = pflag.StringP("config", "c", "", "July config file path.")
)

func main() {
	pflag.Parse()
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}
	log.Compress(viper.GetBool("zap.Compress"))   //是否压缩
	log.MaxSize(viper.GetInt32("zap.MaxSize"))    //单个文件最大体积
	log.Filename(viper.GetString("zap.FileName")) //备份文件名
	log.MaxBackups(viper.GetInt("zap.MaxBackup")) //最大备份留存数
	log.LocalTime(viper.GetBool("zap.LocalTime")) //是否使用本地时间

	logger := log.NewZapLogger()
	e := echo.New()
	Load(e)
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.SugaredLogger.Error("Mongo Break ,maybe you should check it out.")
			}
		}()
		err := store.InitDB().Client().Ping(ctx, readpref.Primary())
		if err != nil {
			panic(err)
		}
	}()
	defer store.Client.Close()
	//e.Use(middleware.CORS())
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.SugaredLogger.Error("Redis down , please test your env.")
			}
		}()
		cacheClient.InitCache()
		_, err := cacheClient.Client.Ping().Result()
		if err != nil {
			panic(err)
		}
	}()
	//go func(e *echo.Echo) {
	//	e.AutoTLSManager.Cache = autocert.DirCache("/conf")
	//	e.Logger.Fatal("TLS service start at port:", e.StartAutoTLS(":443"))
	//}(e)
	e.Logger.Fatal("Service start at port:", e.Start(":2334"))
}
