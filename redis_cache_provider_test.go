package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

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
	Time:     time.Now(),
	UnixTime: UnixTime(time.Now()),
}

func TestRedisCache(t *testing.T) {
	cli := NewRedisCacheProvider(redisC)
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

func TestIncrease(t *testing.T) {
	runtime.GOMAXPROCS(3)
	const DoTimes = 1000
	startNumber := int64(300) //必须非0

	key := "key_" + fmt.Sprint(rand.Int31())
	redisCp.Remove(key) // 确保key 不存在

	wg := sync.WaitGroup{}
	wg.Add(DoTimes)

	successCount := int64(0)
	for i := 0; i < DoTimes; i++ {
		if i == int(startNumber) {
			redisCp.Set(key, startNumber, time.Minute) // 延迟设置键
		}
		go func(t time.Duration) {
			time.Sleep(time.Duration(t) * time.Nanosecond) //延迟执行
			defer wg.Done()
			_, err := redisCp.Increase(key)
			if err == nil {
				atomic.AddInt64(&successCount, 1)
			}
		}(time.Duration(i))
	}
	wg.Wait()
	kv, _ := redisCp.Increase(key) // 测返回值
	if kv-startNumber-1 != successCount {
		t.Errorf("kv: %d  !=   successCount: %d", kv, successCount)
	}
	fmt.Println(kv)
}

func TestIncreaseOrCreate(t *testing.T) {
	key := "key_" + fmt.Sprint(rand.Int31())
	redisCp.Remove(key)

	redisCp.IncreaseOrCreate(key, 10, 1*time.Minute) //1 min
	redisCp.IncreaseOrCreate(key, 10, 2*time.Minute) //1 min

	cmd := redisCp.client.TTL(context.Background(), key)
	if cmd.Val() > 1*time.Minute {
		t.Errorf("expire time error")
	}

	cmd2 := redisCp.client.Get(context.Background(), key)
	if v, _ := cmd2.Int(); v != 20 {
		t.Errorf("expire time error")
	}
}

func TestNewRedisCacheProvider(t *testing.T) {
	type args struct {
		cli redis.Cmdable
	}
	tests := []struct {
		name    string
		args    args
		want    *RedisCacheProvider
		wantErr bool
	}{
		{"Client", args{&redis.Client{}}, &RedisCacheProvider{&redis.Client{}}, false},
		{"Ring client", args{&redis.Ring{}}, &RedisCacheProvider{&redis.Ring{}}, false},
		{"ClusterClient client", args{&redis.ClusterClient{}}, &RedisCacheProvider{&redis.ClusterClient{}}, false},
		{"Tx client", args{&redis.Tx{}}, &RedisCacheProvider{&redis.Tx{}}, false},
		{"nil value", args{nil}, nil, true},
		// 哨兵客户端，不支持。
		//{"client", args{&redis.SentinelClient{}}, &RedisCacheProvider{&redis.Client{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := func() (p *RedisCacheProvider, e error) {
				defer func() {
					err := recover()
					if err != nil {
						e = err.(error)
						p = nil
					}
				}()
				return NewRedisCacheProvider(tt.args.cli), nil
			}()

			if !reflect.DeepEqual(got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("NewRedisCacheProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
