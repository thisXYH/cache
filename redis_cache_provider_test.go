package cache

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

// 测试准备。
func RedisCacheProviderPrepare() {
	p := getNewEveryTime()
	intv := 1
	intvPtr := &intv
	var intf interface{} = intv
	type args struct {
		key   string
		value any
		t     time.Duration
	}
	keys := []args{
		{"bool_true", true, NoExpiration},
		{"bool_false", false, NoExpiration},

		{"set_int8", int8(-8), NoExpiration},
		{"set_int16", int16(-16), NoExpiration},
		{"set_int32", int16(-32), NoExpiration},
		{"set_int64", int64(-64), NoExpiration},
		{"set_int", int(100), NoExpiration},

		{"set_uint8", uint8(8), NoExpiration},
		{"set_uint16", uint16(16), NoExpiration},
		{"set_uint32", uint16(32), NoExpiration},
		{"set_uint64", uint64(64), NoExpiration},
		{"set_uint", uint(100), NoExpiration},

		{"set_float32", float32(123.5), NoExpiration},
		{"set_float64", float64(123.5), NoExpiration},

		{"set_array_int", [3]int{1, 2, 3}, NoExpiration},
		{"set_map_int_string", map[int]string{1: "a", 2: "b", 3: "c"}, NoExpiration},
		{"set_slice_string", []string{"a", "b", "c"}, NoExpiration},
		{"set_struct", data, NoExpiration},

		{"set_ptr", intvPtr, NoExpiration},
		{"set_nil", nil, NoExpiration},
		{"set_interface", intf, NoExpiration},

		{"set_time", tn, NoExpiration},
		{"set_unixTime", UnixTime(tn), NoExpiration},
	}

	for _, v := range keys {
		p.Set(v.key, v.value, v.t)
	}
}

func RedisCacheProviderClearn() {
	p := getNewEveryTime()
	keys := []string{
		"bool_true",
		"bool_false",

		"set_int8",
		"set_int16",
		"set_int32",
		"set_int64",
		"set_int",

		"set_uint8",
		"set_uint16",
		"set_uint32",
		"set_uint64",
		"set_uint",

		"set_float32",
		"set_float64",

		"set_array_int",
		"set_map_int_string",
		"set_slice_string",
		"set_struct",

		"set_ptr",
		"set_nil",
		"set_interface",

		"set_time",
		"set_unixTime",

		"bool_true_create",
		"negative_expireTime",
		"unsupport_kind_complex",
		"unsupport_kind_channel",
		"bool_true_Remove",
		"not_number",
	}

	for _, v := range keys {
		p.Remove(v)
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

func TestRedisCacheProvider_Get(t *testing.T) {
	RedisCacheProviderPrepare()
	defer RedisCacheProviderClearn()

	p := getNewEveryTime()

	get_bool := true
	get_int := 0
	get_uint := uint(0)
	get_float := 1.1
	var get_array [3]int
	get_map := make(map[int]string)
	var get_slice []string
	var get_struct Data
	var get_time time.Time
	var get_unixTime UnixTime

	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name    string
		cli     *RedisCacheProvider
		args    args
		wantErr bool
		wantV   any
	}{
		{"get_bool_true", p, args{"bool_true", &get_bool}, false, true},
		{"get_bool_false", p, args{"bool_false", &get_bool}, false, false},

		{"get_int8", p, args{"set_int8", &get_int}, false, -8},
		{"get_int16", p, args{"set_int16", &get_int}, false, -16},
		{"get_int32", p, args{"set_int32", &get_int}, false, -32},
		{"get_int64", p, args{"set_int64", &get_int}, false, -64},
		{"get_int", p, args{"set_int", &get_int}, false, 100},

		{"get_uint8", p, args{"set_uint8", &get_uint}, false, uint(8)},
		{"get_uint16", p, args{"set_uint16", &get_uint}, false, uint(16)},
		{"get_uint32", p, args{"set_uint32", &get_uint}, false, uint(32)},
		{"get_uint64", p, args{"set_uint64", &get_uint}, false, uint(64)},
		{"get_uint", p, args{"set_uint", &get_uint}, false, uint(100)},

		{"get_float32", p, args{"set_float32", &get_float}, false, 123.5},
		{"get_float64", p, args{"set_float64", &get_float}, false, 123.5},

		{"get_array_int", p, args{"set_array_int", &get_array}, false, [3]int{1, 2, 3}},
		{"get_map_int_string", p, args{"set_map_int_string", &get_map}, false, map[int]string{1: "a", 2: "b", 3: "c"}},
		{"get_slice_string", p, args{"set_slice_string", &get_slice}, false, []string{"a", "b", "c"}},
		{"get_struct", p, args{"set_struct", &get_struct}, false, data},

		{"get_ptr", p, args{"set_ptr", &get_int}, false, 1},

		// 存进去一个 nil，用 任意 ptr 取
		// 只有一层的话应该是取不出来了， 毕竟 (*Data)(nil) 接收不了数据
		// 如果是第二层的估计是可以
		//{"get_nil", p, args{"set_nil", (*Data)(nil)}, false, nil},
		{"get_interface", p, args{"set_interface", &get_int}, false, 1}, //存进去一个 int 的interface , 用 *int 取

		{"get_time", p, args{"set_time", &get_time}, false, tn},
		{"get_unixTime", p, args{"set_unixTime", &get_unixTime}, false, UnixTime(tn)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cli.Get(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisCacheProvider.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.wantV, reflect.ValueOf(tt.args.value).Elem().Interface()) {
				t.Errorf("RedisCacheProvider.Get() = %v, want %v", reflect.ValueOf(tt.args.value).Elem().Interface(), tt.wantV)
			}
		})
	}
}

func TestRedisCacheProvider_MustGet(t *testing.T) {
	RedisCacheProviderPrepare()
	defer RedisCacheProviderClearn()

	p := getNewEveryTime()
	get_bool := true

	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name string
		cli  *RedisCacheProvider
		args args
	}{
		{"empty_key", p, args{"", &get_bool}},
		{"empty_ptr", p, args{"bool_true", (*bool)(nil)}},
		{"not_ptr", p, args{"bool_false", get_bool}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (err error) {
				defer func() {
					e := recover()
					if e != nil {
						err = e.(error)
					}
				}()
				tt.cli.MustGet(tt.args.key, tt.args.value)
				return nil
			}()

			if err == nil {
				t.Errorf(tt.name)
			}
		})
	}
}

func TestRedisCacheProvider_TryGet(t *testing.T) {
	RedisCacheProviderPrepare()
	defer RedisCacheProviderClearn()

	p := getNewEveryTime()
	get_bool := true
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name    string
		cli     *RedisCacheProvider
		args    args
		want    bool
		wantErr bool
	}{
		{"get_exists_key", p, args{"bool_true", &get_bool}, true, false},
		{"get_not_exists_key", p, args{"bool_true_not", &get_bool}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.TryGet(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisCacheProvider.TryGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisCacheProvider.TryGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisCacheProvider_Create(t *testing.T) {
	RedisCacheProviderPrepare()
	defer RedisCacheProviderClearn()

	p := getNewEveryTime()
	type args struct {
		key   string
		value any
		t     time.Duration
	}
	tests := []struct {
		name    string
		cli     *RedisCacheProvider
		args    args
		want    bool
		wantErr bool
	}{
		{"create_exists_key", p, args{"bool_true", true, NoExpiration}, false, false},
		{"create_not_exists_key", p, args{"bool_true_create", true, 1}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.Create(tt.args.key, tt.args.value, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisCacheProvider.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisCacheProvider.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisCacheProvider_MustCreate(t *testing.T) {
	RedisCacheProviderPrepare()
	defer RedisCacheProviderClearn()

	p := getNewEveryTime()
	type args struct {
		key   string
		value any
		t     time.Duration
	}
	tests := []struct {
		name string
		cli  *RedisCacheProvider
		args args
	}{
		{"empty_key", p, args{"", "empty_key", NoExpiration}},
		{"negative_expireTime", p, args{"negative_expireTime", "negative_expireTime", -1}},
		{"unsupport_kind_complex", p, args{"unsupport_kind_complex", complex(1, 2), NoExpiration}},
		{"unsupport_kind_channel", p, args{"unsupport_kind_channel", make(chan int), NoExpiration}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (e error) {
				defer func() {
					err := recover()
					if err != nil {
						e = err.(error)
					}
				}()

				tt.cli.MustCreate(tt.args.key, tt.args.value, tt.args.t)
				return nil
			}()

			if err == nil {
				t.Errorf(tt.name)
			}
		})
	}
}

func TestRedisCacheProvider_Set(t *testing.T) {
	p := getNewEveryTime()
	intv := 1
	intvPtr := &intv
	var intf interface{} = intv

	type args struct {
		key   string
		value any
		t     time.Duration
	}
	tests := []struct {
		name    string
		cli     *RedisCacheProvider
		args    args
		wantErr bool
	}{
		{"set_bool_true", p, args{"bool_true", true, NoExpiration}, false},
		{"set_bool_false", p, args{"bool_false", false, NoExpiration}, false},

		{"set_int8", p, args{"set_int8", int8(-8), NoExpiration}, false},
		{"set_int16", p, args{"set_int16", int16(-16), NoExpiration}, false},
		{"set_int32", p, args{"set_int32", int16(-32), NoExpiration}, false},
		{"set_int64", p, args{"set_int64", int64(-64), NoExpiration}, false},
		{"set_int", p, args{"set_int", int(100), NoExpiration}, false},

		{"set_uint8", p, args{"set_uint8", uint8(8), NoExpiration}, false},
		{"set_uint16", p, args{"set_uint16", uint16(16), NoExpiration}, false},
		{"set_uint32", p, args{"set_uint32", uint16(32), NoExpiration}, false},
		{"set_uint64", p, args{"set_uint64", uint64(64), NoExpiration}, false},
		{"set_uint", p, args{"set_uint", uint(100), NoExpiration}, false},

		{"set_float32", p, args{"set_float32", float32(123.5), NoExpiration}, false},
		{"set_float64", p, args{"set_float64", float64(123.5), NoExpiration}, false},

		{"set_array_int", p, args{"set_array_int", [3]int{1, 2, 3}, NoExpiration}, false},
		{"set_map_int_string", p, args{"set_map_int_string", map[int]string{1: "a", 2: "b", 3: "c"}, NoExpiration}, false},
		{"set_slice_string", p, args{"set_slice_string", []string{"a", "b", "c"}, NoExpiration}, false},
		{"set_struct", p, args{"set_struct", data, NoExpiration}, false},

		{"set_ptr", p, args{"set_ptr", intvPtr, NoExpiration}, false},
		{"set_nil", p, args{"set_nil", nil, NoExpiration}, false},
		{"set_interface", p, args{"set_interface", intf, NoExpiration}, false},

		{"set_time", p, args{"set_time", tn, NoExpiration}, false},
		{"set_unixTime", p, args{"set_unixTime", UnixTime(tn), NoExpiration}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cli.Set(tt.args.key, tt.args.value, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("RedisCacheProvider.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	RedisCacheProviderClearn()
}

func TestRedisCacheProvider_MustSet(t *testing.T) {
	p := getNewEveryTime()

	type args struct {
		key   string
		value any
		t     time.Duration
	}
	tests := []struct {
		name      string
		cli       *RedisCacheProvider
		args      args
		wantPanic bool
	}{
		{"empty_key", p, args{"", "empty_key", NoExpiration}, true},
		{"negative_expireTime", p, args{"negative_expireTime", "negative_expireTime", -1}, true},
		{"unsupport_kind_complex", p, args{"unsupport_kind_complex", complex(1, 2), NoExpiration}, true},
		{"unsupport_kind_channel", p, args{"unsupport_kind_channel", make(chan int), NoExpiration}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (e error) {
				defer func() {
					err := recover()
					if err != nil {
						e = err.(error)
					}
				}()
				tt.cli.MustSet(tt.args.key, tt.args.value, tt.args.t)
				return nil
			}()

			log.Println(tt.name, ":", err)

			if tt.wantPanic && err == nil {
				t.Fail()
			}
		})
	}
}

func TestRedisCacheProvider_Remove(t *testing.T) {
	RedisCacheProviderPrepare()
	defer RedisCacheProviderClearn()

	p := getNewEveryTime()
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		cli     *RedisCacheProvider
		args    args
		want    bool
		wantErr bool
	}{
		{"remov_exists_key", p, args{"bool_true"}, true, false},
		{"remov_exists_key", p, args{"bool_false"}, true, false},

		{"remov_exists_key", p, args{"set_int8"}, true, false},
		{"remov_exists_key", p, args{"set_int16"}, true, false},
		{"remov_exists_key", p, args{"set_int32"}, true, false},
		{"remov_exists_key", p, args{"set_int64"}, true, false},
		{"remov_exists_key", p, args{"set_int"}, true, false},

		{"remov_exists_key", p, args{"set_uint8"}, true, false},
		{"remov_exists_key", p, args{"set_uint16"}, true, false},
		{"remov_exists_key", p, args{"set_uint32"}, true, false},
		{"remov_exists_key", p, args{"set_uint64"}, true, false},
		{"remov_exists_key", p, args{"set_uint"}, true, false},

		{"remov_exists_key", p, args{"set_float32"}, true, false},
		{"remov_exists_key", p, args{"set_float64"}, true, false},

		{"remov_exists_key", p, args{"set_array_int"}, true, false},
		{"remov_exists_key", p, args{"set_map_int_string"}, true, false},
		{"remov_exists_key", p, args{"set_slice_string"}, true, false},
		{"remov_exists_key", p, args{"set_struct"}, true, false},

		{"remov_exists_key", p, args{"set_ptr"}, true, false},
		{"remov_exists_key", p, args{"set_nil"}, true, false},
		{"remov_exists_key", p, args{"set_interface"}, true, false},

		{"remov_exists_key", p, args{"set_time"}, true, false},
		{"remov_exists_key", p, args{"set_unixTime"}, true, false},

		{"remov_not_exists_key", p, args{"bool_true_Remove"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.Remove(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisCacheProvider.Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisCacheProvider.Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisCacheProvider_MustRemove(t *testing.T) {
	RedisCacheProviderPrepare()
	defer RedisCacheProviderClearn()

	p := getNewEveryTime()
	type args struct {
		key string
	}
	tests := []struct {
		name string
		cli  *RedisCacheProvider
		args args
	}{
		{"empty_key", p, args{""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (e error) {
				defer func() {
					err := recover()
					if err != nil {
						e = err.(error)
					}
				}()
				tt.cli.MustRemove(tt.args.key)
				return nil
			}()
			if err == nil {
				t.Errorf(tt.name)
			}
		})
	}
}

// TestRedisCacheProvider_Increase
// 因为是通过事务实现的，主要测试并发，事务的完整性
// 主要点：程序中的成功 必须等于 redis中的成功。
func TestRedisCacheProvider_Increase(t *testing.T) {
	p := getNewEveryTime()
	runtime.GOMAXPROCS(3)
	const DoTimes = 1000
	startNumber := int64(300) // 起点，必须非0

	key := "key_" + fmt.Sprint(rand.Int31())
	p.Remove(key) // 确保key 不存在

	wg := sync.WaitGroup{}
	wg.Add(DoTimes)

	successCount := int64(0)
	for i := 0; i < DoTimes; i++ {
		if i == int(startNumber) {
			p.Set(key, startNumber, time.Minute) // 延迟设置键
		}
		go func(t time.Duration) {
			time.Sleep(time.Duration(t)) //延迟执行
			defer wg.Done()
			_, err := p.Increase(key)
			if err == nil {
				atomic.AddInt64(&successCount, 1)
			}
		}(time.Duration(i))
	}
	wg.Wait()
	kv, _ := p.Increase(key) // 测返回值
	if kv-startNumber-1 != successCount {
		t.Errorf("kv: %d  !=   successCount: %d", kv, successCount)
	}
}

func TestRedisCacheProvider_MustIncrease(t *testing.T) {
	p := getNewEveryTime()
	p.Set("not_number", "not number", 10*time.Second)

	type args struct {
		key string
	}
	tests := []struct {
		name string
		cli  *RedisCacheProvider
		args args
	}{
		{"not_number", p, args{"not_number"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := func() (e error) {
				defer func() {
					err := recover()
					if err != nil {
						e = err.(error)
					}
				}()
				tt.cli.MustIncrease(tt.args.key)
				return nil
			}()

			if err == nil {
				t.Errorf(tt.name)
			}
		})
	}
}

func TestRedisCacheProvider_IncreaseOrCreate(t *testing.T) {
	p := getNewEveryTime()
	key := "key_" + fmt.Sprint(rand.Int31())
	p.Remove(key)

	p.IncreaseOrCreate(key, 10, 1*time.Minute) //1 min
	p.IncreaseOrCreate(key, 10, 2*time.Minute) //1 min

	cmd := p.client.TTL(context.Background(), key)
	if cmd.Val() > 1*time.Minute {
		t.Errorf("IncreaseOrCreate: expire time error")
	}

	after := 0
	p.MustGet(key, &after)
	if after != 20 {
		t.Errorf("IncreaseOrCreate: number error, %d", after)
	}
}

func TestRedisCacheProvider_MustIncreaseOrCreate(t *testing.T) {
	p := getNewEveryTime()
	p.Set("not_number", "not number", 10*time.Second)

	type args struct {
		key       string
		increment int64
		t         time.Duration
	}
	tests := []struct {
		name string
		cli  *RedisCacheProvider
		args args
	}{
		{"not_number", p, args{"not_number", 10, 1 * time.Minute}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (e error) {
				defer func() {
					err := recover()
					if err != nil {
						e = err.(error)
					}
				}()
				tt.cli.MustIncreaseOrCreate(tt.args.key, tt.args.increment, tt.args.t)
				return nil
			}()

			if err == nil {
				t.Errorf(tt.name)
			}
		})
	}
}
