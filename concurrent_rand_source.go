package cache

import (
	"math/rand"
	"sync"
)

// concurrentRandSource 线程安全的 rand.Source 实现。
//  issue: https://github.com/golang/go/issues/3611 .
type concurrentRandSource struct {
	rand.Source
	m sync.Mutex
}

var (
	_ rand.Source = (*concurrentRandSource)(nil)
)

func (c *concurrentRandSource) Int63() int64 {
	c.m.Lock()
	defer c.m.Unlock()

	return c.Source.Int63()
}

func (c *concurrentRandSource) Seed(seed int64) {
	c.m.Lock()
	defer c.m.Unlock()

	c.Source.Seed(seed)
}

func newConcurrentRandSource(seed int64) rand.Source {
	return &concurrentRandSource{
		Source: rand.NewSource(seed),
		m:      sync.Mutex{},
	}
}
