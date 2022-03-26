package cache

import (
	"fmt"
	"reflect"
	"time"
)

// Level2CacheProvider 实现简单的两级缓存，不支持Increase， 两个层次的缓存使用相同的缓存key,
// 所以两个层级的缓存需要使用不同的缓存提供器，防止相互覆盖。
//
// 当一级获取不到，将从二级获取（一般来说，一级回收间隔更短），
// 并且根据二级缓存更新一级缓存。
//
// 实际上可以看成以 Level 2 为主，Level 1 为辅助，提高访问性能。
type Level2CacheProvider struct {
	level1 CacheProvider // 一级缓存
	level2 CacheProvider // 二级缓存

	// expireTime 一级缓存的过期时间，
	// 其他 Api 的缓存时间为二级缓存的缓存时间，
	// 要求二级缓存的过期时间 > 一级缓存。
	expireTime *Expiration
}

// NewLevel2CacheProvider 新建一个二级缓存提供器。
func NewLevel2CacheProvider(l1, l2 CacheProvider, expireTime *Expiration) *Level2CacheProvider {
	if expireTime == nil {
		panic(fmt.Errorf("'expireTime' must not be nil"))
	}

	return &Level2CacheProvider{l1, l2, expireTime}
}

// implement CacheProvider.Get .
func (p *Level2CacheProvider) Get(key string, value any) error {
	_, err := p.TryGet(key, value)
	return err
}

// implement CacheProvider.TryGet .
func (p *Level2CacheProvider) TryGet(key string, value any) (result bool, err error) {
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

// implement CacheProvider.Create .
func (p *Level2CacheProvider) Create(key string, value any, t time.Duration) (result bool, err error) {
	// 异常或者二级缓存 key 存在。
	if result, err = p.level2.Create(key, value, t); err != nil || !result {
		return
	}

	// 二级缓存 key，已经不存在了，就连带覆盖一级缓存。
	p.setLevel1(key, value)
	return true, nil
}

// implement CacheProvider.Set .
func (p *Level2CacheProvider) Set(key string, value any, t time.Duration) error {
	err := p.level2.Set(key, value, t)
	if err != nil {
		return err
	}

	p.setLevel1(key, value)
	return err
}

// implement CacheProvider.Remove .
func (p *Level2CacheProvider) Remove(key string) (bool, error) {
	result, err := p.level2.Remove(key)

	p.level1.Remove(key)

	return result, err
}

// implement CacheProvider.Increase, not supported, will panic!
func (p *Level2CacheProvider) Increase(key string) (int64, error) {
	panic("not supported")
}

// implement CacheProvider.IncreaseOrCreate, not supported, will panic!
func (p *Level2CacheProvider) IncreaseOrCreate(key string, increment int64, t time.Duration) (int64, error) {
	panic("not supported")
}

// setLevel1 is set cache for Level 1.
func (p *Level2CacheProvider) setLevel1(key string, value any) error {
	return p.level1.Set(key, value, p.expireTime.NextExpireTime())
}
