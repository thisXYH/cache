package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestMemoryCacheProvider(t *testing.T) {
	cache := NewMemoryCacheProvider(10 * time.Second)

	// cache.Set("int", int64(10), time.Minute)

	// var i int8
	// cache.MustGet("int", &i)
	// fmt.Println(i)

	var m map[string]int = make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	cache.MustSet("map", m, time.Minute)

	var m2 map[string]int
	cache.MustGet("map", &m2)

	fmt.Println(m2)
}
