package cache

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/cmstar/go-conv"
	c "github.com/patrickmn/go-cache"
)

// MemoryCacheProvider 内存类型的缓存提供器。
type MemoryCacheProvider struct {
	cache *c.Cache // 线程安全的缓存
	mu    sync.RWMutex
}

// NewMemoryCacheProvider.
func NewMemoryCacheProvider(cleanupInterval time.Duration) *MemoryCacheProvider {
	// 限制清理周期 >= 1 sec 防止负载过高，以及锁缓存。
	if cleanupInterval < time.Second {
		panic(fmt.Errorf("'cleanupInterval' must be greater than 1 second"))
	}
	return &MemoryCacheProvider{c.New(cleanupInterval, cleanupInterval), sync.RWMutex{}}
}

var (
	_ CacheProvider = (*MemoryCacheProvider)(nil)
)

// implement CacheProvider.Get .
func (cp *MemoryCacheProvider) Get(key string, value any) error {
	_, err := cp.TryGet(key, value)

	return err
}

// implement CacheProvider.MustGet .
func (cp *MemoryCacheProvider) MustGet(key string, value any) {
	err := cp.Get(key, value)
	if err != nil {
		panic(err)
	}
}

// implement CacheProvider.TryGet .
func (cp *MemoryCacheProvider) TryGet(key string, value any) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key must not be empty")
	}

	cp.mu.RLock()
	defer cp.mu.RUnlock()

	item, exists := cp.cache.Get(key)
	if !exists {
		return false, nil
	}
	itemT := reflect.TypeOf(item)

	// 基础类型使用转换。
	if conv.IsPrimitiveKind(itemT.Kind()) {
		err := conv.Convert(item, value)
		return true, err
	}

	// 非基础类型，直接设置值， 反射不能设置 unexposed field。
	reflect.ValueOf(value).Elem().Set(reflect.ValueOf(item))
	return true, nil
}

// implement CacheProvider.Create .
func (cp *MemoryCacheProvider) Create(key string, value any, t time.Duration) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key must not be empty")
	}
	cp.mu.Lock()
	defer cp.mu.Unlock()

	t = cp.legalExpireTime(t)
	err := cp.cache.Add(key, value, t)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// implement CacheProvider.MustCreate .
func (cp *MemoryCacheProvider) MustCreate(key string, value any, t time.Duration) bool {
	v, err := cp.Create(key, value, t)
	if err != nil {
		panic(err)
	}

	return v
}

// implement CacheProvider.Set .
func (cp *MemoryCacheProvider) Set(key string, value any, t time.Duration) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}
	cp.mu.Lock()
	defer cp.mu.Unlock()

	t = cp.legalExpireTime(t)
	cp.cache.Set(key, value, t)
	return nil
}

// implement CacheProvider.MustSet .
func (cp *MemoryCacheProvider) MustSet(key string, value any, t time.Duration) {
	err := cp.Set(key, value, t)
	if err != nil {
		panic(err)
	}
}

// implement CacheProvider.Remove .
func (cp *MemoryCacheProvider) Remove(key string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key must not be empty")
	}
	cp.mu.Lock()
	defer cp.mu.Unlock()

	_, exists := cp.cache.Get(key)
	cp.cache.Delete(key)
	return exists, nil
}

// implement CacheProvider.MustRemove .
func (cp *MemoryCacheProvider) MustRemove(key string) bool {
	v, err := cp.Remove(key)
	if err != nil {
		panic(err)
	}
	return v
}

// implement CacheProvider.Increase .
func (cp *MemoryCacheProvider) Increase(key string) (int64, error) {
	if key == "" {
		return 0, fmt.Errorf("key must not be empty")
	}
	cp.mu.Lock()
	defer cp.mu.Unlock()

	v, expireTime, found := cp.cache.GetWithExpiration(key)
	if !found {
		return 0, fmt.Errorf("cache key does not exist: %s", key)
	}

	if _, ok := v.(int64); ok {
		return cp.cache.IncrementInt64(key, 1)
	}

	var v64 int64
	switch v := v.(type) {
	case int:
		v64 = int64(v)
	case int8:
		v64 = int64(v)
	case int16:
		v64 = int64(v)
	case int32:
		v64 = int64(v)
	case int64:
		v64 = int64(v)
	default:
		return 0, fmt.Errorf("unsupport type to increase: %s", reflect.TypeOf(v).Kind())
	}

	// 更新 key 的数据类型，并且避免过期时间重置。
	cp.cache.Set(key, v64, time.Duration(expireTime.UnixNano()-time.Now().UnixNano()))

	return cp.cache.IncrementInt64(key, 1)
}

// implement CacheProvider.MustIncrease .
func (cp *MemoryCacheProvider) MustIncrease(key string) int64 {
	v, err := cp.Increase(key)
	if err != nil {
		panic(err)
	}
	return v
}

// implement CacheProvider.IncreaseOrCreate .
func (cp *MemoryCacheProvider) IncreaseOrCreate(key string, increment int64, t time.Duration) (int64, error) {
	if key == "" {
		return 0, fmt.Errorf("key must not be empty")
	}
	cp.mu.Lock()
	defer cp.mu.Unlock()

	v, expireTime, found := cp.cache.GetWithExpiration(key)
	if !found {
		cp.cache.Set(key, increment, t)
		return increment, nil
	}

	if _, ok := v.(int64); ok {
		return cp.cache.IncrementInt64(key, increment)
	}

	var v64 int64
	switch v := v.(type) {
	case int:
		v64 = int64(v)
	case int8:
		v64 = int64(v)
	case int16:
		v64 = int64(v)
	case int32:
		v64 = int64(v)
	case int64:
		v64 = int64(v)
	default:
		return 0, fmt.Errorf("unsupport type to increase: %s", reflect.TypeOf(v).Kind())
	}

	// 更新 key 的数据类型，并且避免过期时间重置。
	r := v64 + increment
	cp.cache.Set(key, r, time.Duration(expireTime.UnixNano()-time.Now().UnixNano()))
	return r, nil
}

// implement CacheProvider.MustIncreaseOrCreate .
func (cp *MemoryCacheProvider) MustIncreaseOrCreate(key string, increment int64, t time.Duration) int64 {
	v, err := cp.IncreaseOrCreate(key, increment, t)
	if err != nil {
		panic(err)
	}
	return v
}

func (*MemoryCacheProvider) legalExpireTime(t time.Duration) time.Duration {
	if t < 0 {
		panic(fmt.Errorf("expire time must not be less than 0"))
	}

	if t == NoExpiration {
		return c.NoExpiration
	}

	return t
}
