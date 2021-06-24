package cache

import (
	"testing"
	"time"
)

var cacheOp = NewCacheOperation("go", "cache:test", 2, redisCp, 5*time.Minute, 1*time.Minute)

func TestKeyOp(t *testing.T) {

	value := "test valueaa"
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
