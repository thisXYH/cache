package cache_test

import (
	"github.com/thisXYH/cache"
	"testing"
)

var cacheProvider = cache.NewRedisCacheProvider(RedisClient)

var cacheOp = cache.NewCacheOperation("go", "cache:test", 2, cacheProvider, 0)

func TestKeyOp(t *testing.T) {
	value := "test value"
	keyOp := cacheOp.Key("keyOp", "haha")
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
