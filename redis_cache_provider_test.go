package cache_test

import (
	"encoding/json"
	"fmt"
	"github.com/thisXYH/cache"
	"testing"
	"time"
)

var cli = cache.NewRedisCacheProvider(RedisClient)

type Person struct {
	Name string
	Age  int
}

type Data struct {
	Bool bool

	Int   int
	Int8  int8
	Int16 int16
	Int32 int32
	Int64 int64

	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Uintptr uintptr

	Float32 float32
	Float64 float64

	Array [2]int
	Map   map[string]int

	Slice  []int
	String string

	Person Person

	Time     time.Time
	UnixTime cache.UnixTime
}

var data = Data{
	Bool:    true,
	Int:     1,
	Int8:    2,
	Int16:   3,
	Int32:   4,
	Int64:   5,
	Uint:    6,
	Uint8:   7,
	Uint16:  8,
	Uint32:  9,
	Uint64:  10,
	Uintptr: 11,
	Float32: 12.12,
	Float64: 13.13,

	Array:    [2]int{16, 16},
	Map:      map[string]int{"A": 17, "B": 18},
	Slice:    []int{19, 20},
	String:   "21",
	Person:   Person{"Jerry", 22},
	Time:     time.Now(),
	UnixTime: cache.UnixTime(time.Now()),
}

func TestRedisCache(t *testing.T) {
	dv, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(dv))

	// 确保key不存在
	key := "go.redis.test.cli"
	cli.MustRemove(key)

	// 获取不存在的缓存
	var v Data = Data{Bool: true} //设置一个非 0值
	if cli.MustGet(key, &v); !v.Bool {
		t.Error("在没有缓存key的情况下，接收者值被改变")
	}

	// 尝试获取不存在的缓存
	if exist, err := cli.TryGet(key, &v); exist || err != nil || !v.Bool {
		t.Error("预期之外的行为", exist, v.Bool, err)
	}

	// 设置缓存
	cli.MustSet(key, data, 0)

	// 获取缓存值，对比是否和之前的一样。
	cli.MustGet(key, &v)
	tv, err := json.Marshal(v)
	if err != nil {
		t.Error(err)
	}

	if string(dv) != string(tv) {
		t.Error("存进去、取出来 值发生了改变")
	}

	// SETNX
	result, err := cli.Create(key, Data{Int: 99}, 0)
	if err != nil {
		t.Error(err)
	}

	if result {
		t.Error("预期 key 存在，添加缓存失败")
	}

	cli.MustGet(key, &v)
	if v.Int == 99 {
		t.Error("预期 key 存在，添加缓存失败，缓存值还是之前的值")
	}

	// 删除缓存key， 存在和不存在时的不同返回值。
	r1 := cli.MustRemove(key)
	r2 := cli.MustRemove(key)
	if !r1 || r2 {
		t.Error("预期 r1 == true , r2 == false")
	}
}
