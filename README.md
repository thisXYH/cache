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


