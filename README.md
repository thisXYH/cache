# cache
提供一套缓存操作的语义，以及一些已实现的缓存提供器。

## 功能
* [X] 缓存key管理
* [X] 随机过期时间管理
* [X] 自定义缓存提供器
* 缓存提供器
    * [X] Redis缓存
    * [X] memory缓存
    * [X] 二级缓存, 不支持Increase。

## 快速开始
```bash
go get -u github.com/thisXYH/cache@latest
```
## 注意事项
* redis cahce provider 的测试用例需要自己提供redis连接信息才能测试成功。
* concurrent rand source 的测试用例是并发测试，可能需要跑多次才能成功。

## 缓存语义
```go
// CacheProvider 提供一套缓存语义。
type CacheProvider interface {
	// Get 获取指定缓存值。
	//  @key: cache key.
	//  @value: receive value.
	// return: 如果key存在，value被更新成对应值, 反之value值不做改变。
	Get(key string, value any) error

	// MustGet 是 Get 的 panic 版。
	MustGet(key string, value any)

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

	// MustCreate 是 Create 的 panic 版。
	MustCreate(key string, value any, t time.Duration) bool

	// Set 设置或者更新缓存。
	//  @key: cache key.
	//  @value: receive value.
	//  @t: 过期时长， 0表不过期。
	Set(key string, value any, t time.Duration) error

	// MustSet 是 Set 的 panic 版。
	MustSet(key string, value any, t time.Duration)

	// Remove 移除指定缓存,
	//  @key: cache key.
	// return: true 成功移除；false 缓存不存在。
	Remove(key string) (bool, error)

	// MustRemove 是 Remove 的 panic 版。
	MustRemove(key string) bool

	// Increase 为已存在的指定缓存的值（必须是整数）增加1。
	//  @key: cache key.
	// return: 符合条件返回增加后的值，反之返回默认值，以及对应的 error。
	Increase(key string) (int64, error)

	// MustIncrease 是 Increase 的 panic 版。
	MustIncrease(key string) int64

	// Increase 为指定缓存的值增加一个增量(负数==减法)，如果不存在则创建该缓存。
	//  @key: cache key.
	//  @increment: 增量，如果 key 不存在，则当成新缓存的 value。
	//  @t: 过期时长， 0表不过期。
	// return: 返回增加后的值。
	IncreaseOrCreate(key string, increment int64, t time.Duration) (int64, error)

	// MustIncreaseOrCreate 是 IncreaseOrCreate 的 panic 版。
	MustIncreaseOrCreate(key string, increment int64, t time.Duration) int64
}
```

## 使用示例
> [base use sample](https://github.com/thisXYH/cache/blob/main/internal/sample_test.go)
````go
package cache_example

import (
	"fmt"
	"github.com/thisXYH/cache"
	"time"
)

// 创建一个缓存操作对象。
//  指定key的组成: go:cache:test<_flag1><_flag2>
//  指定缓存提供器：内存缓存，设置一分钟清理一次过期 key。
//  指定缓存的过期时间为，10分钟 再加上 2分钟的随机量 上下波动。
//
// 直接声明成全局变量。
var cacheOp *cache.Operation = cache.NewOperation(
	"go", "cache:test", 2,
	cache.NewMemoryCacheProvider(1*time.Minute),
	cache.NewExpirationFromMinute(10, 2))

func doSomething() {
	// 获取 key 操作对象，这个对象中包含了组装好的完整缓存key。
	// cacheOp 指定了两个flag，这边就必须传两个参数，多、少都不行。
	// go:cache:test_unixTime_1 .
	key := cacheOp.Key(time.Now(), true)

	// 获取完整缓存key
	fmt.Println(key.Key)
	// output: go:cache:test_1625123485000_1

	// 设置缓存。
	key.Set("hello world!")

	// 获取缓存。
	var value string
	key.Get(&value)

	// Get 的 panic版本，如果key不存在，直接panic。
	key.MustGet(&value)

	fmt.Println(value)
	// output: hello world!

	// 支持基础类型互转，存进去一个 int8  用 int 去接。
	// 转换的实现逻辑可参考：github.com/cmstar/go-conv 。
	key.Set(int8(8))
	var intV int
	key.Get(&intV)

	fmt.Println(intV)
	// output: 8
}
````