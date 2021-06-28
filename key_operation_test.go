package cache

import (
	"fmt"
	"time"
)

func ExampleKeyOperation() {
	// 创建一个缓存操作对象。
	//  指定key的组成: go:cache:test<_flag1><_flag2>
	//  指定缓存提供器：内存缓存，切一分钟清理一次过期key
	//  指定缓存的过期时间为，10分钟 再加上 2分钟 的随机量 上下波动
	//
	// 正常情况下该变量声明成全局变量，不需要重复声明。
	cacheOp := NewOperation("go", "cache:test", 2, NewMemoryCacheProvider(1*time.Minute), NewExpireTimeFromMinute(10, 2))

	// 获取 key 操作对象
	// 这个对象中包含了组装好的完整换缓存key，
	// go:cache:test_unixTime_1，
	// cacheOp 指定了两个flag，这边就必须传两个参数，多、少都不行。
	key := cacheOp.Key(time.Now(), true)

	// 获取完整缓存key
	// fmt.Println(key.Key)

	key.Set("hellow word!")

	var value string
	key.Get(&value)

	fmt.Println(value)
	// output: hellow word!
}
