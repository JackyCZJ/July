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
	//addrSlice := viper.GetStringSlice("redis.cluster")
	//Ring = redis.NewRing(&redis.RingOptions{
	//	Addrs: addrMap,
	//	DB:    0,
	//})
	opt := redis.Options{
		Addr: "redis1:6379",
	}
	client := redis.NewClient(&opt)
	//opt := redis.ClusterOptions{}
	//addrSlice := viper.GetStringSlice("redis.cluster")
	//opt.RouteRandomly = true
	//opt.ClusterSlots = func() ([]redis.ClusterSlot, error) {
	//	slots := []redis.ClusterSlot{{
	//		Start: 0,
	//		End:   4999,
	//		Nodes: []redis.ClusterNode{{
	//			ID:   "",
	//			Addr: addrSlice[0],
	//		}, {
	//			ID:   "",
	//			Addr: addrSlice[3],
	//		},
	//		},
	//	}, {
	//		Start: 5000,
	//		End:   9999,
	//		Nodes: []redis.ClusterNode{{
	//			ID:   "",
	//			Addr: addrSlice[1],
	//		}, {
	//			ID:   "",
	//			Addr: addrSlice[4],
	//		}},
	//	}, {
	//		Start: 10000,
	//		End:   16383,
	//		Nodes: []redis.ClusterNode{{
	//			ID:   "",
	//			Addr: addrSlice[2],
	//		}, {
	//			ID:   "",
	//			Addr: addrSlice[5],
	//		}},
	//	}}
	//	return slots, nil
	//}
	//Cluster = redis.NewClusterClient(&opt)
	//ring := redis.NewRing(&redis.RingOptions{
	//	Addrs: viper.GetStringMapString("redis.docker_cluster"),
	//})
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
