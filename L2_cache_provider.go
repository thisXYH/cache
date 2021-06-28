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
type L2CacheProvider struct {
	level1 CacheProvider // 一级缓存
	level2 CacheProvider // 二级缓存

	// expireTime 一级缓存的过期时间，
	// 其他Api的缓存时间为二级缓存的缓存时间，
	// 要求二级缓存的过期时间 > 一级缓存。
	expireTime *ExpireTime
}

func NewL2CacheProvider(l1, l2 CacheProvider, expireTime *ExpireTime) *L2CacheProvider {
	if expireTime == nil {
		panic(fmt.Errorf("'expireTime' must not be nil"))
	}

	return &L2CacheProvider{l1, l2, expireTime}
}

// Get 获取指定缓存值。
//  @key: cache key
//  @value: receive value
//
//  return: 如果key存在，value被更新成对应值, 反之value值不做改变。
func (p *L2CacheProvider) Get(key string, value any) error {
	_, err := p.TryGet(key, value)
	return err
}

// MustGet 是 Get 的 panic 版。
func (p *L2CacheProvider) MustGet(key string, value any) {
	err := p.Get(key, value)
	if err != nil {
		panic(err)
	}
}

// TryGet 尝试获取指定缓存。
//  @key: cache key
//  @value: receive value
//
// return: 若key存在，value被更新成对应值，返回true，
// 反之value值不做改变，返回false。
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

// Create 仅当缓存键不存在时，创建缓存。
//  @key: cache key
//  @value: receive value
//  @t: 过期时长， 0表不过期
// return: true表示创建了缓存；false说明缓存已经存在了。
func (p *L2CacheProvider) Create(key string, value any, t time.Duration) (result bool, err error) {
	// 异常或者二级缓存key还存在
	if result, err = p.level2.Create(key, value, t); err != nil || !result {
		return
	}

	// 二级缓存key，已经不存在了，就连带覆盖一级缓存。
	p.setLevel1(key, value)
	return true, nil
}

// MustCreate 是 Create 的 panic 版。
func (p *L2CacheProvider) MustCreate(key string, value any, t time.Duration) bool {
	result, err := p.Create(key, value, t)
	if err != nil {
		panic(err)
	}
	return result
}

// Set 设置或者更新缓存，
//  @key: cache key
//  @value: receive value
//  @t: 过期时长， 0表不过期
func (p *L2CacheProvider) Set(key string, value any, t time.Duration) error {
	err := p.level2.Set(key, value, t)
	if err != nil {
		return err
	}

	p.setLevel1(key, value)
	return err
}

// MustSet 是 Set 的 panic 版。
func (p *L2CacheProvider) MustSet(key string, value any, t time.Duration) {
	err := p.Set(key, value, t)
	if err != nil {
		panic(err)
	}
}

// Remove 移除指定缓存,
//  @key: cache key
// return: true成功移除，false缓存不存在。
func (p *L2CacheProvider) Remove(key string) (bool, error) {
	reslut, err := p.level2.Remove(key)

	p.level1.Remove(key)

	return reslut, err
}

// MustRemove 是 Remove 的 panic 版。
func (p *L2CacheProvider) MustRemove(key string) bool {
	v, err := p.Remove(key)
	if err != nil {
		panic(err)
	}
	return v
}

// Increase 为已存在的指定缓存的值（必须是整数）增加1。
//  @key: cache key
// return: 符合条件返回增加后的值，反之返回默认值，以及对应的error.
func (p *L2CacheProvider) Increase(key string) (int64, error) {
	panic("not implemented") // TODO: Implement
}

// MustIncrease 是 Increase 的 panic 版。
func (p *L2CacheProvider) MustIncrease(key string) int64 {
	panic("not implemented") // TODO: Implement
}

// Increase 为指定缓存的值增加一个增量(负数==减法)，如果不存在则创建该缓存。
//  @key: cache key
//  @increment: 增量，如果key不存在，则当成新缓存的value
//  @t: 过期时长， 0表不过期
// return: 返回增加后的值。
func (p *L2CacheProvider) IncreaseOrCreate(key string, increment int64, t time.Duration) (int64, error) {
	panic("not implemented") // TODO: Implement
}

// MustIncreaseOrCreate 是 IncreaseOrCreate 的 panic 版。
func (p *L2CacheProvider) MustIncreaseOrCreate(key string, increment int64, t time.Duration) int64 {
	panic("not implemented") // TODO: Implement
}

// setLevel1 设置一级缓存。
func (p *L2CacheProvider) setLevel1(key string, value any) error {
	return p.level1.Set(key, value, p.expireTime.NextExpireTime())
}
