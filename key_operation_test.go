package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestKeyOperation(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)
	ns := "ns"
	prefix := "prefix"

	t.Run("2-flag", func(t *testing.T) {
		op := NewOperation(ns, prefix, 2, provider, CacheExpirationZero)

		key := op.Key("a", 1)
		fmt.Println(key.Key)

		var res time.Time
		if key.MustTryGet(&res) {
			t.Fatal("key should not be")
		}

		v := time.Date(2022, 03, 27, 18, 55, 0, 0, time.UTC)
		key.MustSet(v)

		if !key.MustTryGet(&res) {
			t.Fatal("key should be")
		}

		if res != v {
			t.Fatal("value err")
		}

		v2 := time.Now()
		if key.MustCreate(v2) {
			t.Fatal("key exist, cannot set cache")
		}

		if !key.MustRemove() {
			t.Fatal("remove fail")
		}
	})
}

func TestKeyOperationT(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)
	ns := "ns"
	prefix := "prefix"

	t.Run("string-int", func(t *testing.T) {
		op := NewOperation2[string, int, time.Time](ns, prefix, provider, CacheExpirationZero)

		key := op.Key("a", 1)
		fmt.Println(key.Key)

		res, ok := key.MustTryGet()
		if ok {
			t.Fatal("key should not be")
		}

		v := time.Date(2022, 03, 27, 18, 55, 0, 0, time.UTC)
		key.MustSet(v)

		res, ok = key.MustTryGet()
		if !ok {
			t.Fatal("key should be")
		}

		if res != v {
			t.Fatal("value err")
		}

		v2 := time.Now()
		if key.MustCreate(v2) {
			t.Fatal("key exist, cannot set cache")
		}

		if !key.MustRemove() {
			t.Fatal("remove fail")
		}
	})
}
