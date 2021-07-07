package cache

import (
	"fmt"
	"reflect"
	"time"
)

// L2CacheProvider 实现简单的两级缓存，不支持Increase， 两个层次的缓存使用相同的缓存key，
// 所以两个层级的缓存需要使用不同的缓存提供器，防止相互覆盖。
//
// 当一级获取不到，将从二级获取（一般来说，一级回收间隔更短），
// 并且根据二级缓存更新一级缓存。
//
// 实际上可以看成以 Level 2 为主，Level 1 为辅助，提高访问性能。
type L2CacheProvider struct {
	level1 CacheProvider // 一级缓存
	level2 CacheProvider // 二级缓存

	// expireTime 一级缓存的过期时间，
	// 其他Api的缓存时间为二级缓存的缓存时间，
	// 要求二级缓存的过期时间 > 一级缓存。
	expireTime *Expiration
}

func NewL2CacheProvider(l1, l2 CacheProvider, expireTime *Expiration) *L2CacheProvider {
	if expireTime == nil {
		panic(fmt.Errorf("'expireTime' must not be nil"))
	}

	return &L2CacheProvider{l1, l2, expireTime}
}

// implement ICacheProvider.Get
func (p *L2CacheProvider) Get(key string, value any) error {
	_, err := p.TryGet(key, value)
	return err
}

// implement ICacheProvider.MustGet
func (p *L2CacheProvider) MustGet(key string, value any) {
	err := p.Get(key, value)
	if err != nil {
		panic(err)
	}
}

// implement ICacheProvider.TryGet
func (p *L2CacheProvider) TryGet(key string, value any) (result bool, err error) {
	if result, err = p.level1.TryGet(key, value); err != nil || result {
		return
	}

	if result, err = p.level2.TryGet(key, value); err != nil {
		return
	}

	if result {
		// value 一定是指针。
		p.setLevel1(key, reflect.ValueOf(value).Elem().Interface())
		return
	}

	return false, nil
}

// implement ICacheProvider.Create
func (p *L2CacheProvider) Create(key string, value any, t time.Duration) (result bool, err error) {
	// 异常或者二级缓存key存在
	if result, err = p.level2.Create(key, value, t); err != nil || !result {
		return
	}

	// 二级缓存key，已经不存在了，就连带覆盖一级缓存。
	p.setLevel1(key, value)
	return true, nil
}

// implement ICacheProvider.MustCreate
func (p *L2CacheProvider) MustCreate(key string, value any, t time.Duration) bool {
	result, err := p.Create(key, value, t)
	if err != nil {
		panic(err)
	}
	return result
}

// implement ICacheProvider.Set
func (p *L2CacheProvider) Set(key string, value any, t time.Duration) error {
	err := p.level2.Set(key, value, t)
	if err != nil {
		return err
	}

	p.setLevel1(key, value)
	return err
}

// implement ICacheProvider.MustSet
func (p *L2CacheProvider) MustSet(key string, value any, t time.Duration) {
	err := p.Set(key, value, t)
	if err != nil {
		panic(err)
	}
}

// implement ICacheProvider.Remove
func (p *L2CacheProvider) Remove(key string) (bool, error) {
	reslut, err := p.level2.Remove(key)

	p.level1.Remove(key)

	return reslut, err
}

// implement ICacheProvider.MustRemove
func (p *L2CacheProvider) MustRemove(key string) bool {
	v, err := p.Remove(key)
	if err != nil {
		panic(err)
	}
	return v
}

// implement ICacheProvider.Increase,
// not implemented, will panic
func (p *L2CacheProvider) Increase(key string) (int64, error) {
	panic("not implemented") // TODO: Implement
}

// implement ICacheProvider.MustIncrease,
// not implemented, will panic
func (p *L2CacheProvider) MustIncrease(key string) int64 {
	panic("not implemented") // TODO: Implement
}

// implement ICacheProvider.IncreaseOrCreate,
// not implemented, will panic
func (p *L2CacheProvider) IncreaseOrCreate(key string, increment int64, t time.Duration) (int64, error) {
	panic("not implemented") // TODO: Implement
}

// implement ICacheProvider.MustIncreaseOrCreate,
// not implemented, will panic
func (p *L2CacheProvider) MustIncreaseOrCreate(key string, increment int64, t time.Duration) int64 {
	panic("not implemented") // TODO: Implement
}

// setLevel1 is set cache for Level 1。
func (p *L2CacheProvider) setLevel1(key string, value any) error {
	return p.level1.Set(key, value, p.expireTime.NextExpireTime())
}
