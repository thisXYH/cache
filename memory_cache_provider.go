package cache

import (
	"cache/conv"
	"fmt"
	"reflect"
	"sync"
	"time"

	c "github.com/patrickmn/go-cache"
)

// 内存 类型的缓存提供器,
// 数据的组织方式，基础类型直接使用
type MemoryCacheProvider struct {
	cache *c.Cache // 线程安全的缓存
	mu    sync.RWMutex
}

// NewMemoryCacheProvider
func NewMemoryCacheProvider(cleanupInterval time.Duration) *MemoryCacheProvider {
	// 限制清理周期 >= 1 sec 防止负载过高，以及锁缓存。
	if cleanupInterval < time.Second {
		panic(fmt.Errorf("'cleanupInterval' must be greater than 1 second"))
	}
	return &MemoryCacheProvider{c.New(cleanupInterval, cleanupInterval), sync.RWMutex{}}
}

var (
	_ ICacheProvider = (*MemoryCacheProvider)(nil)
)

// implement ICacheProvider.Get
func (cp *MemoryCacheProvider) Get(key string, value any) error {
	_, err := cp.TryGet(key, value)

	return err
}

// implement ICacheProvider.MustGet
func (cp *MemoryCacheProvider) MustGet(key string, value any) {
	err := cp.Get(key, value)
	if err != nil {
		panic(err)
	}
}

// implement ICacheProvider.TryGet
func (cp *MemoryCacheProvider) TryGet(key string, value any) (succ bool, err error) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	v, exists := cp.cache.Get(key)
	if !exists {
		return false, nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr {
		return false, fmt.Errorf("'value' is not pointer")
	}

	if !rv.IsValid() {
		return false, fmt.Errorf("'value' is nil")
	}

	rv = rv.Elem()
	var temp interface{}
	switch rv.Kind() {
	case reflect.Bool, reflect.String,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		temp, err = conv.Convert(v, rv.Type())
		if err != nil {
			return false, nil
		}
	default:
		if !rv.CanSet() {
			return false, fmt.Errorf("%t can't set value", rv.Type())
		}

		temp = v
	}

	rv.Set(reflect.ValueOf(temp))
	return true, nil
}

// implement ICacheProvider.Create
func (cp *MemoryCacheProvider) Create(key string, value any, t time.Duration) (bool, error) {
	cp.mu.Lock()
	defer cp.mu.Lock()

	t = cp.legalExpireTime(t)
	err := cp.cache.Add(key, value, t)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// implement ICacheProvider.MustCreate
func (cp *MemoryCacheProvider) MustCreate(key string, value any, t time.Duration) bool {
	v, err := cp.Create(key, value, t)
	if err != nil {
		panic(err)
	}

	return v
}

// implement ICacheProvider.Set
func (cp *MemoryCacheProvider) Set(key string, value any, t time.Duration) error {
	cp.mu.Lock()
	defer cp.mu.Lock()

	t = cp.legalExpireTime(t)
	cp.cache.Set(key, value, t)
	return nil
}

// implement ICacheProvider.MustSet
func (cp *MemoryCacheProvider) MustSet(key string, value any, t time.Duration) {
	err := cp.Set(key, value, t)
	if err != nil {
		panic(err)
	}
}

// implement ICacheProvider.Remove
func (cp *MemoryCacheProvider) Remove(key string) (bool, error) {
	cp.mu.Lock()
	defer cp.mu.Lock()

	_, exists := cp.cache.Get(key)
	cp.cache.Delete(key)
	return exists, nil
}

// implement ICacheProvider.MustRemove
func (cp *MemoryCacheProvider) MustRemove(key string) bool {
	v, err := cp.Remove(key)
	if err != nil {
		panic(err)
	}
	return v
}

// implement ICacheProvider.Increase
func (cp *MemoryCacheProvider) Increase(key string) (int64, error) {
	cp.mu.Lock()
	defer cp.mu.Lock()

	return cp.cache.IncrementInt64(key, 1)
}

// implement ICacheProvider.MustIncrease
func (cp *MemoryCacheProvider) MustIncrease(key string) int64 {
	v, err := cp.Increase(key)
	if err != nil {
		panic(err)
	}
	return v
}

// implement ICacheProvider.IncreaseOrCreate
func (cp *MemoryCacheProvider) IncreaseOrCreate(key string, increment int64, t time.Duration) (int64, error) {
	cp.mu.Lock()
	defer cp.mu.Lock()

	_, exists := cp.cache.Get(key)
	if exists {
		return cp.cache.IncrementInt64(key, increment)
	}

	cp.cache.Set(key, increment, t)
	return increment, nil
}

// implement ICacheProvider.MustIncreaseOrCreate
func (cp *MemoryCacheProvider) MustIncreaseOrCreate(key string, increment int64, t time.Duration) int64 {
	v, err := cp.IncreaseOrCreate(key, increment, t)
	if err != nil {
		panic(err)
	}
	return v
}

func (*MemoryCacheProvider) legalExpireTime(t time.Duration) time.Duration {
	if t < 0 {
		panic(fmt.Errorf("expire time must not be letter than 0"))
	}

	if t == NoExpiration {
		return c.NoExpiration
	}

	return t
}
