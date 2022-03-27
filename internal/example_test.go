package cache_example

import (
	"fmt"
	"time"

	"github.com/thisXYH/cache"
)

func ExampleKeyOperation() {
	// 初始化一个 内存缓存提供器。
	provider := cache.NewMemoryCacheProvider(time.Second)

	// 初始化一个 缓存操作对象, 通常这个对象会被用作全局, 不需要每次使用都创建。
	op := cache.NewOperation("ns", "prefix", 2, provider, cache.CacheExpirationZero)

	// 根据指定的数量给定 unique flag, 获取对应的缓存 key 操作对象。
	key := op.Key("a", 1)
	fmt.Println(key.Key)

	var res time.Time
	if key.MustTryGet(&res) {
		// 如果 key 存在，进入这个代码块(当前一定不存在，展示作用)。
		panic("key should not be")
	}

	v := time.Date(2022, 03, 27, 18, 55, 0, 0, time.UTC)
	key.MustSet(v)

	if !key.MustTryGet(&res) {
		panic("key should be")
	}

	if res != v {
		panic("value err")
	}

	v2 := time.Now()

	// 当 key 不存在的时候，设置缓存返回 true, 反之 false 。
	if key.MustCreate(v2) {
		panic("key exist, cannot set cache")
	}

	// 最后删除缓存, 删除成功返回 true, 反之 false 。
	if !key.MustRemove() {
		panic("remove fail")
	}

	// output:
	// ns:prefix_a_1
}

func ExampleKeyOperationT() {
	// 初始化一个 内存缓存提供器。
	provider := cache.NewMemoryCacheProvider(time.Second)

	// 初始化一个 缓存操作对象, 通常这个对象会被用作全局, 不需要每次使用都创建。
	op := cache.NewOperation2[string, int, time.Time]("ns", "prefix", provider, cache.CacheExpirationZero)

	// 根据指定的数量给定 unique flag, 获取对应的缓存 key 操作对象。
	key := op.Key("a", 1)
	fmt.Println(key.Key)

	res, ok := key.MustTryGet()
	if ok {
		// 如果 key 存在，进入这个代码块(当前一定不存在，展示作用)。
		panic("key should not be")
	}

	v := time.Date(2022, 03, 27, 18, 55, 0, 0, time.UTC)
	key.MustSet(v)

	res, ok = key.MustTryGet()
	if !ok {
		panic("key should be")
	}

	if res != v {
		panic("value err")
	}

	v2 := time.Now()

	// 当 key 不存在的时候，设置缓存返回 true, 反之 false 。
	if key.MustCreate(v2) {
		panic("key exist, cannot set cache")
	}

	// 最后删除缓存, 删除成功返回 true, 反之 false 。
	if !key.MustRemove() {
		panic("remove fail")
	}

	// output:
	// ns:prefix_a_1
}
