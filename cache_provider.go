package cache

import (
	"time"
)

type any interface{}

const (
	//不过期。
	NoExpiration time.Duration = 0
)

// ICacheProvider 提供一套缓存语义。
type ICacheProvider interface {
	// Get 获取指定缓存值,
	// 如果key存在，value被更新成对应值，
	// 反之value值不做改变。
	Get(key string, value any) error

	// MustGet 是 Get 的 panic 版。
	MustGet(key string, value any)

	// TryGet 尝试获取指定缓存，
	// 若key存在，value被更新成对应值，返回true，
	// 反之value值不做改变，返回false。
	TryGet(key string, value any) (bool, error)

	// Create 仅当缓存键不存在时，创建缓存，
	// t 过期时长， 0 表不过期。
	// return: true表示创建了缓存；false说明缓存已经存在了。
	Create(key string, value any, t time.Duration) (bool, error)

	// MustCreate 是 Create 的 panic 版。
	MustCreate(key string, value any, t time.Duration) bool

	// Set 设置或者更新缓存，
	// t 过期时长， 0 表不过期。
	Set(key string, value any, t time.Duration) error

	// MustSet 是 Set 的 panic 版。
	MustSet(key string, value any, t time.Duration)

	// Remove 移除指定缓存,
	// return: true成功移除，false缓存不存在。
	Remove(key string) (bool, error)

	// MustRemove 是 Remove 的 panic 版。
	MustRemove(key string) bool
}
