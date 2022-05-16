package cache

import (
	"strings"
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
		if key.Key != "ns:prefix_a_1" {
			t.Fatal("error cache key be generated")
		}

		var timeDefaultValue time.Time
		var res time.Time

		key.MustGet(&res)
		if res != timeDefaultValue {
			t.Fatal("when key no exit, value cannot be modified")
		}

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

	t.Run("increase", func(t *testing.T) {
		op := NewOperation(ns, prefix, 1, provider, CacheExpirationZero)
		key := op.Key("increase")

		_, err := key.Increase()
		if !strings.HasPrefix(err.Error(), "cache key does not exist") {
			t.Fatal("increase fail:", err.Error())
		}

		if key.MustIncreaseOrCreate(2) != 2 {
			t.Fatal("increase or create fail")
		}

		if key.MustIncrease() != 3 {
			t.Fatal("increase fail")
		}

		key.MustRemove()
	})
}

func TestKeyOperationT(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)
	ns := "ns"
	prefix := "prefix"

	t.Run("string-int", func(t *testing.T) {
		op := NewOperation2[string, int, time.Time](ns, prefix, provider, CacheExpirationZero)

		key := op.Key("a", 1)
		if key.Key != "ns:prefix_a_1" {
			t.Fatal("error cache key be generated")
		}

		var timeDefaultValue time.Time
		var res time.Time

		res = key.MustGet()
		if res != timeDefaultValue {
			t.Fatal("when key no exit, value should be default")
		}

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

	t.Run("increaseT", func(t *testing.T) {
		op := NewOperation1[string, int](ns, prefix, provider, CacheExpirationZero)
		key := op.Key("increaseT")

		_, err := key.Increase()
		if !strings.HasPrefix(err.Error(), "cache key does not exist") {
			t.Fatal("increase fail:", err.Error())
		}

		if key.MustIncreaseOrCreate(2) != 2 {
			t.Fatal("increase or create fail")
		}

		if key.MustIncrease() != 3 {
			t.Fatal("increase fail")
		}

		key.MustRemove()
	})
}
