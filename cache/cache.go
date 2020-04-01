package cache

import (
	"time"

	"github.com/jackyczj/July/log"

	"github.com/go-redis/redis/v7"

	"github.com/go-redis/cache/v7"

	"github.com/vmihailenco/msgpack/v4"
)

var (
	//Ring *redis.Ring
	// global cache
	cc *cache.Codec

	Cluster *redis.ClusterClient

	Client *redis.Client
)

// 初始化缓存
func InitCache() {
	opt := redis.Options{
		Addr: "redis1:6379",
	}
	client := redis.NewClient(&opt)
	Client = client
	cc = &cache.Codec{
		Redis: client,
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
		log.Logworker.Error(err.Error())
	}
}

// getCc 读缓存
func GetCc(key string, pointer interface{}) error {
	return cc.Get(key, pointer)
}

// delCc 清缓存
func DelCc(key string) {
	if err := cc.Delete(key); err != nil {
		log.Logworker.Error(err.Error())
	}
}

// cleanCc 批量清除一类缓存
func CleanCc(cate string) {
	if cate == "" {
		log.Logworker.Error("someone try to clean all cache keys")
		return
	}
	i := 0
	for _, key := range Cluster.Keys(cate + "*").Val() {
		DelCc(key)
		i++
	}
	log.Logworker.Infof("delete %d %s cache", i, cate)
}
