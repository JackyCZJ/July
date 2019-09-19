package cache

import (
	"net/http"
	"time"

	"github.com/go-redis/redis"

	"github.com/go-redis/cache"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	"github.com/vmihailenco/msgpack"
)

var (
	// redis client
	rdb *redis.Client
	// global cache
	cc *cache.Codec
)

// 初始化缓存
func initCache() {
	cc = &cache.Codec{
		Redis: rdb,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

// setCc 写缓存
func setCc(key string, object interface{}, exp time.Duration) {
	cc.Set(&cache.Item{
		Key:        key,
		Object:     object,
		Expiration: exp,
	})
}

// getCc 读缓存
func getCc(key string, pointer interface{}) error {
	return cc.Get(key, pointer)
}

// delCc 清缓存
func delCc(key string) {
	cc.Delete(key)
}

// cleanCc 批量清除一类缓存
func cleanCc(cate string) {
	if cate == "" {
		logrus.Error("someone try to clean all cache keys")
		return
	}
	i := 0
	for _, key := range rdb.Keys(cate + "*").Val() {
		delCc(key)
		i++
	}
	logrus.Infof("delete %d %s cache", i, cate)
}

func deleteCache(c echo.Context) error {
	cate := c.Param("cate")
	switch cate {
	case "token":
		cleanCc("token")
	case "all":
		cleanCc("token")
	default:
		return echo.NewHTTPError(400, "InvalidID", "请在URL中提供合法的缓存类型")
	}
	return c.NoContent(http.StatusNoContent)
}
