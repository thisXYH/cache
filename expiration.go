package cache

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const NoExpiration time.Duration = 0                      // 缓存不过期。
var CacheExpirationZero *Expiration = NewExpiration(0, 0) // 缓存不过期。

// Expiration 缓存过期时间。
type Expiration struct {
	// expireTime 过期时长。
	baseExpireTime time.Duration

	// 过期时间随机量。
	randomRangeTime time.Duration

	// 加锁处理线程安全问题。
	// issue: https://github.com/golang/go/issues/3611 .
	rand      *rand.Rand
	randMutex sync.Mutex
}

// NewExpiration 新建缓存过期时间。
//  @baseExpireTime: 基准过期时长，0表不过期。
//  @randomRangeTime: 随机过期市场，0表不随机，否则baseExpireTime将增加[-randomRangeTime, +randomRangeTime]。
func NewExpiration(baseExpireTime, randomRangeTime time.Duration) *Expiration {
	if baseExpireTime < 0 {
		panic(fmt.Errorf("'baseExpireTime' must not be less than 0"))
	}

	if randomRangeTime < 0 {
		panic(fmt.Errorf("'randomRangeTime' must not be less than 0"))
	}

	if baseExpireTime < randomRangeTime {
		panic(fmt.Errorf("'baseExpireTime' must not be less than 'randomRangeTime'"))
	}

	var randTemp *rand.Rand
	if randomRangeTime != 0 {
		randTemp = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	return &Expiration{baseExpireTime, randomRangeTime, randTemp, sync.Mutex{}}
}

// NewExpirationFromMillisecond 以毫秒为单位创建缓存过期时间。
func NewExpirationFromMillisecond(baseExpireTime, randomRangeTime int64) *Expiration {
	return NewExpiration(time.Duration(baseExpireTime)*time.Millisecond, time.Duration(randomRangeTime)*time.Millisecond)
}

// NewExpirationFromSecond 以秒为单位创建缓存过期时间。
func NewExpirationFromSecond(baseExpireTime, randomRangeTime int64) *Expiration {
	return NewExpiration(time.Duration(baseExpireTime)*time.Second, time.Duration(randomRangeTime)*time.Second)
}

// NewExpirationFromMinute 以分钟为单位创建缓存过期时间。
func NewExpirationFromMinute(baseExpireTime, randomRangeTime int64) *Expiration {
	return NewExpiration(time.Duration(baseExpireTime)*time.Minute, time.Duration(randomRangeTime)*time.Minute)
}

// NewExpirationFromHour 以小时为单位创建缓存过期时间。
func NewExpirationFromHour(baseExpireTime, randomRangeTime int64) *Expiration {
	return NewExpiration(time.Duration(baseExpireTime)*time.Hour, time.Duration(randomRangeTime)*time.Hour)
}

// NextExpireTime 获取一个新的过期时间，如果存在随机量的话，返回值已经过随机量计算。
func (c *Expiration) NextExpireTime() time.Duration {
	if c.baseExpireTime == NoExpiration {
		return NoExpiration
	}

	if c.rand == nil {
		return c.baseExpireTime
	}

	c.randMutex.Lock()
	defer c.randMutex.Unlock()

	randomRangeTimeInt := int64(c.randomRangeTime)
	if c.rand.Int31n(2) == 0 {
		return c.baseExpireTime - time.Duration(c.rand.Int63n(randomRangeTimeInt))
	} else {
		return c.baseExpireTime + time.Duration(c.rand.Int63n(randomRangeTimeInt))
	}
}
