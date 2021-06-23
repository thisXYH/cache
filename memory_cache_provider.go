package cache

import (
	"cache/conv"
	"errors"
	"reflect"
	"sync"
	"time"
)

// 内存 类型的缓存提供器,
// 数据的组织方式，基础类型直接使用
type MemoryCacheProvider struct {
	data           map[string]memoryCacheData
	mu             sync.Mutex
	isClearning    bool
	thresholdCount int
	count          int
}

type memoryCacheData struct {
	value any
	// 过期时间。
	expireTime int64
}

func NewMemoryCacheProvider(count int) *MemoryCacheProvider {
	return &MemoryCacheProvider{
		data:           make(map[string]memoryCacheData, 1024),
		thresholdCount: 1024,
		count:          1024,
		isClearning:    false,
		mu:             sync.Mutex{},
	}
}

var (
	_ ICacheProvider = (*MemoryCacheProvider)(nil)
)

func (cp *MemoryCacheProvider) Get(key string, value any) error {
	_, err := cp.TryGet(key, value)
	return err
}

func (cp *MemoryCacheProvider) MustGet(key string, value any) {
	err := cp.Get(key, value)
	if err != nil {
		panic(err)
	}
}

func (cp *MemoryCacheProvider) TryGet(key string, value any) (bool, error) {
	v, ok := cp.data[key]
	if !ok {
		return false, nil
	}

	// 过期
	if cp.expireIfNeeded(key, v) {
		return false, nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr {
		return false, errors.New("MemoryCacheProvider: Unmarshal(non-pointer " + rv.Type().String() + ")")
	}

	if !rv.IsValid() {
		return false, errors.New("MemoryCacheProvider: Unmarshal(nil " + rv.Type().String() + ")")
	}
	rv = rv.Elem()
	temp, err := conv.Convert(v.value, rv.Type())
	if err != nil {
		return false, nil
	}

	rv.Set(reflect.ValueOf(temp))
	return true, nil
}

func (cp *MemoryCacheProvider) Create(key string, value any, t time.Duration) (bool, error) {
	v, ok := cp.data[key]
	if ok && !cp.expireIfNeeded(key, v) {
		return false, nil
	}

	cp.data[key] = memoryCacheData{value, time.Now().UnixNano() + int64(t)}
	return true, nil
}

func (cp *MemoryCacheProvider) MustCreate(key string, value any, t time.Duration) bool {
	v, err := cp.Create(key, value, t)
	if err != nil {
		panic(err)
	}

	return v
}

func (cp *MemoryCacheProvider) Set(key string, value any, t time.Duration) error {
	cp.data[key] = memoryCacheData{value, time.Now().UnixNano() + int64(t)}
	return nil
}

func (cp *MemoryCacheProvider) MustSet(key string, value any, t time.Duration) {
	err := cp.Set(key, value, t)
	if err != nil {
		panic(err)
	}
}

func (cp *MemoryCacheProvider) Remove(key string) (bool, error) {
	_, ok := cp.data[key]
	if ok {
		delete(cp.data, key)
	}
	return true, nil
}

func (cp *MemoryCacheProvider) MustRemove(key string) bool {
	v, err := cp.Remove(key)
	if err != nil {
		panic(err)
	}
	return v
}

// expireIfNeeded 过期缓存如果需要的话。
// true 过期
// false 未过期
func (cp *MemoryCacheProvider) expireIfNeeded(key string, v memoryCacheData) bool {
	if v.expireTime <= time.Now().UnixNano() {
		delete(cp.data, key)
		return true
	}
	return false
}
