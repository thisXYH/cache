package cache

import (
	"testing"
)

var cacheOp = NewOperation("go", "cache:test", 2, redisCp, NewExpireTimeFromMinute(5, 2))

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
