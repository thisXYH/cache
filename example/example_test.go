package example

import (
	"cache"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

// 实例化一个缓存提供器。
var cacheProvider cache.ICacheProvider = cache.NewRedisCacheProvider(redis.NewClient(&redis.Options{
	Addr:     "HOST:PORT",
	Password: "PASSWORD", // no password set
	DB:       0,          // use default DB
}))

// 实例化一个缓存操作对象。
var cacheOp = cache.NewCacheOperation("go", "cache:test", 3, cacheProvider, 5*time.Minute, 1*time.Minute)

func TestExample(t *testing.T) {
	value := "test value"
	keyOp := cacheOp.Key("keyOp", time.Now(), true)
	fmt.Println("cache key:", keyOp.Key)

	keyOp.MustSet(value)

	var recv string
	keyOp.MustGet(&recv)

	if recv != value {
		t.Error(recv)
	}

	keyOp.MustRemove()

	exist, _ := keyOp.TryGet(&recv)
	if exist {
		t.Error(exist)
	}
}
