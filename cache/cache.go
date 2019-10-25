package cache

import (
	"net/http"
	"time"

	"github.com/spf13/viper"

	"github.com/go-redis/redis"

	"github.com/go-redis/cache"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/vmihailenco/msgpack"
)

var (
	//Ring *redis.Ring
	// global cache
	cc *cache.Codec

	Cluster *redis.ClusterClient
)

// 初始化缓存
func InitCache() {
	//addrSlice := viper.GetStringSlice("redis.cluster")
	//Ring = redis.NewRing(&redis.RingOptions{
	//	Addrs: addrMap,
	//	DB:    0,
	//})
	opt := redis.ClusterOptions{}
	addrSlice := viper.GetStringSlice("redis.cluster")
	opt.ClusterSlots = func() ([]redis.ClusterSlot, error) {
		slots := []redis.ClusterSlot{{
			Start: 0,
			End:   4999,
			Nodes: []redis.ClusterNode{{
				Addr: addrSlice[0],
			}},
		}, {
			Start: 5000,
			End:   9999,
			Nodes: []redis.ClusterNode{{
				Addr: addrSlice[1],
			}},
		}, {
			Start: 10000,
			End:   16383,
			Nodes: []redis.ClusterNode{{
				Addr: addrSlice[2],
			}},
		}}
		return slots, nil
	}
	Cluster = redis.NewClusterClient(&opt)

	cc = &cache.Codec{
		Redis: Cluster,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

// setCc 写缓存
func SetCc(key string, object interface{}, exp time.Duration) {
	err := cc.Set(&cache.Item{
		Key:        key,
		Object:     object,
		Expiration: exp,
	})
	if err != nil {
		logrus.Error(err.Error())
	}
}

// getCc 读缓存
func GetCc(key string, pointer interface{}) error {
	return cc.Get(key, pointer)
}

// delCc 清缓存
func DelCc(key string) {
	if err := cc.Delete(key); err != nil {
		logrus.Error(err.Error())
	}
}

// cleanCc 批量清除一类缓存
func CleanCc(cate string) {
	if cate == "" {
		logrus.Error("someone try to clean all cache keys")
		return
	}
	i := 0
	for _, key := range Cluster.Keys(cate + "*").Val() {
		DelCc(key)
		i++
	}
	logrus.Infof("delete %d %s cache", i, cate)
}

func DeleteCache(c echo.Context) error {
	cate := c.Param("cate")
	switch cate {
	case "token":
		CleanCc("token")
	case "all":
		CleanCc("token")
	default:
		return echo.NewHTTPError(400, "InvalidID", "请在URL中提供合法的缓存类型")
	}
	return c.NoContent(http.StatusNoContent)
}
