package cache

import (
	"time"
)

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
