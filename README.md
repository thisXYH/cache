# cache

[![License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/thisXYH/cache.svg)](https://pkg.go.dev/github.com/thisXYH/cache)

提供一套缓存操作的语义，以及一些已实现的缓存提供器。

## 功能
* [X] 缓存key管理
* [X] 随机过期时间管理
* [X] 自定义缓存提供器
* 缓存提供器
    * [X] Redis缓存
    * [X] memory缓存
    * [X] 二级缓存, 不支持Increase
* 支持泛型(version >= v1.1.0)

## 快速开始
```bash
go get -u github.com/thisXYH/cache@latest
```
## 注意事项
* redis cahce provider 的测试用例需要自己提供redis连接信息才能测试成功。通过 `conf_test.go` 的 `getNewEveryTime` 调整 redis 配置，默认使用 `127.0.0.1:6379` 。
* concurrent rand source 的测试用例是并发测试，可能需要跑多次才能成功。
* 新建 redis client 时 `Options.MaxRetries` 用于指定客户端失败重试，当重试次数不等于 1 的时候，由于客户端会自动重试，可能会导致 `Increase` 和 `IncreaseOrCreate` 语义不准。

## 缓存语义
```go
// CacheProvider 提供一套缓存语义。
type CacheProvider interface {
	// Get 获取指定缓存值。
	//  @key: cache key.
	//  @value: receive value.
	// return: 如果key存在，value被更新成对应值, 反之value值不做改变。
	Get(key string, value any) error

	// TryGet 尝试获取指定缓存。
	//  @key: cache key.
	//  @value: receive value.
	// return: 若 key 存在，value 被更新成对应值，返回 true；反之 value 值不做改变，返回 false。
	TryGet(key string, value any) (bool, error)

	// Create 仅当缓存键不存在时，创建缓存。
	//  @key: cache key.
	//  @value: receive value.
	//  @t: 过期时长， 0表不过期。
	// return: true表示创建了缓存；false说明缓存已经存在了。
	Create(key string, value any, t time.Duration) (bool, error)

	// Set 设置或者更新缓存。
	//  @key: cache key.
	//  @value: receive value.
	//  @t: 过期时长， 0表不过期。
	Set(key string, value any, t time.Duration) error

	// Remove 移除指定缓存,
	//  @key: cache key.
	// return: true 成功移除；false 缓存不存在。
	Remove(key string) (bool, error)

	// Increase 为已存在的指定缓存的值（必须是整数）增加1。
	//  @key: cache key.
	// return: 符合条件返回增加后的值，反之返回默认值，以及对应的 error。
	Increase(key string) (int64, error)

	// Increase 为指定缓存的值增加一个增量(负数==减法)，如果不存在则创建该缓存。
	//  @key: cache key.
	//  @increment: 增量，如果 key 不存在，则当成新缓存的 value。
	//  @t: 过期时长， 0表不过期。
	// return: 返回增加后的值。
	IncreaseOrCreate(key string, increment int64, t time.Duration) (int64, error)
}
```

## 使用示例
> [base use sample](https://github.com/thisXYH/cache/blob/main/internal/example_test.go)
````go
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
````