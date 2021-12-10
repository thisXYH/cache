package cache

import (
	"time"

	"github.com/go-redis/redis/v8"
)

func getNewEveryTime() *RedisCacheProvider {
	cli := redis.NewClient(&redis.Options{
		Addr: "redis.cache.tongbu.com:16381",
		//Password: "CDDayDbWNJ6zW3p1weIuv9", // no password set
		DB: 0, // use default DB

		// 最大重试次数，-1 不重试
		// 对于一些非幂等的命令，执行重试是不合理的，比如 incr
		MaxRetries: -1,
	})

	return NewRedisCacheProvider(cli)
}

// 测试的时候使用。
var tn, _ = time.ParseInLocation("2006-01-02 15:04:05", "2006-01-02 15:04:05", time.Local)

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
	UnixTime UnixTime
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
	Time:     tn,
	UnixTime: UnixTime(tn),
}
