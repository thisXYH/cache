package cache

import (
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
	kind  reflect.Kind
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
	panic("not implemented") // TODO: Implement
}

func (cp *MemoryCacheProvider) MustGet(key string, value any) {
	panic("not implemented") // TODO: Implement
}

func (cp *MemoryCacheProvider) TryGet(key string, value any) (bool, error) {
	v, ok := cp.data[key]
	if !ok {
		return false, nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr {
		return false, errors.New("MemoryCacheProvider: Unmarshal(non-pointer " + rv.Type().String() + ")")
	}

	if !rv.IsValid() {
		return false, errors.New("MemoryCacheProvider: Unmarshal(nil " + rv.Type().String() + ")")
	}

	if rv.Kind() != v.kind {
		return false, errors.New("MemoryCacheProvider: type not equal")
	}

	rv.Elem().Set(reflect.ValueOf(v.value))

	return true, nil
}

func (cp *MemoryCacheProvider) Create(key string, value any, t time.Duration) (bool, error) {
	panic("not implemented") // TODO: Implement
}

func (cp *MemoryCacheProvider) MustCreate(key string, value any, t time.Duration) bool {
	panic("not implemented") // TODO: Implement
}

func (cp *MemoryCacheProvider) Set(key string, value any, t time.Duration) error {
	panic("not implemented") // TODO: Implement
}

func (cp *MemoryCacheProvider) MustSet(key string, value any, t time.Duration) {
	panic("not implemented") // TODO: Implement
}

func (cp *MemoryCacheProvider) Remove(key string) (bool, error) {
	panic("not implemented") // TODO: Implement
}

func (cp *MemoryCacheProvider) MustRemove(key string) bool {
	panic("not implemented") // TODO: Implement
}
