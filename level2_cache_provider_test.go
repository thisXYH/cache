package cache

import (
	"testing"
	"time"
)

func TestNewLevel2CacheProvider(t *testing.T) {
	rp := getNewEveryTime()

	t.Run("overview", func(t *testing.T) {
		mp := NewMemoryCacheProvider(time.Second)
		p := NewLevel2CacheProvider(mp, rp, CacheExpirationZero)

		key := "Level2CacheProvider_test_overview"
		var v int

		err := p.Get(key, &v)
		if err != nil {
			t.Fatal(err)
		}
		if v != 0 {
			t.Fatal("when key no exit, value cannot be modified")
		}

		ok, err := p.TryGet(key, &v)
		if ok {
			t.Fatal("key should not be")
		}

		err = p.Set(key, 10, NoExpiration)
		if err != nil {
			t.Fatal(err)
		}

		ok, err = p.TryGet(key, &v)
		if !ok {
			t.Fatal("key should be")
		}

		if v != 10 {
			t.Fatal("value err")
		}

		ok, err = p.Create(key, 11, NoExpiration)
		if ok {
			t.Fatal("key exist, cannot set cache")
		}

		ok, err = p.Remove(key)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("remove fail")
		}

		var v2 int
		p.Get(key, &v2)
		if v2 != 0 {
			t.Fatal("remove fail, key exist")
		}
	})

	t.Run("expiration", func(t *testing.T) {
		mp := NewMemoryCacheProvider(time.Second)

		// 1级缓存 三秒过期。
		p := NewLevel2CacheProvider(mp, rp, NewExpirationFromSecond(3, 0))
		key := "Level2CacheProvider_test_expiration"

		// 二级缓存 一分钟过期。
		p.Set(key, 10, time.Minute)

		// 查看当前 一级缓存中是否有存在。
		var v int
		p.level1.Get(key, &v)
		if v != 10 {
			t.Fatal("key not exit, in level 1 cache")
		}

		// 让一级缓存过期.
		time.Sleep(5 * time.Second)
		v = 0
		p.level1.Get(key, &v)
		if v != 0 {
			t.Fatal("key exit, in level 1 cache")
		}

		// 从二级缓存读取，并且覆盖一级缓存。
		p.Get(key, &v)
		if v != 10 {
			t.Fatal("key not exit, in level 2 cache")
		}

		v = 0
		p.level1.Get(key, &v)
		if v != 10 {
			t.Fatal("key not exit, in level 1 cache")
		}

		p.Remove(key)
	})
}
