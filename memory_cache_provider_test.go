package cache

import (
	"reflect"
	"testing"
	"time"
)

// 测试用
var mp = NewMemoryCacheProvider(2 * time.Second)

func MemoryCacheProviderPrepare() {
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
		mp.Set(v.key, v.value, v.t)
	}
}

func MemoryCacheProviderClearn() {
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
	}

	for _, v := range keys {
		mp.Remove(v)
	}
}

func TestNewMemoryCacheProvider(t *testing.T) {
	type args struct {
		cleanupInterval time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"negative", args{-10}, true},
		{"positive", args{time.Second}, false},
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
				NewMemoryCacheProvider(tt.args.cleanupInterval)
				return nil
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("NewMemoryCacheProvider error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryCacheProvider_Get(t *testing.T) {
	MemoryCacheProviderPrepare()
	defer MemoryCacheProviderClearn()

	get_bool := true
	get_int := 0
	//get_int_ptr := &get_int
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
		cp      *MemoryCacheProvider
		args    args
		wantErr bool
		wantV   any
	}{
		{"get_bool_true", mp, args{"bool_true", &get_bool}, false, true},
		{"get_bool_false", mp, args{"bool_false", &get_bool}, false, false},

		{"get_int8", mp, args{"set_int8", &get_int}, false, -8},
		{"get_int16", mp, args{"set_int16", &get_int}, false, -16},
		{"get_int32", mp, args{"set_int32", &get_int}, false, -32},
		{"get_int64", mp, args{"set_int64", &get_int}, false, -64},
		{"get_int", mp, args{"set_int", &get_int}, false, 100},

		{"get_uint8", mp, args{"set_uint8", &get_uint}, false, uint(8)},
		{"get_uint16", mp, args{"set_uint16", &get_uint}, false, uint(16)},
		{"get_uint32", mp, args{"set_uint32", &get_uint}, false, uint(32)},
		{"get_uint64", mp, args{"set_uint64", &get_uint}, false, uint(64)},
		{"get_uint", mp, args{"set_uint", &get_uint}, false, uint(100)},

		{"get_float32", mp, args{"set_float32", &get_float}, false, 123.5},
		{"get_float64", mp, args{"set_float64", &get_float}, false, 123.5},

		{"get_array_int", mp, args{"set_array_int", &get_array}, false, [3]int{1, 2, 3}},
		{"get_map_int_string", mp, args{"set_map_int_string", &get_map}, false, map[int]string{1: "a", 2: "b", 3: "c"}},
		{"get_slice_string", mp, args{"set_slice_string", &get_slice}, false, []string{"a", "b", "c"}},
		{"get_struct", mp, args{"set_struct", &get_struct}, false, data},

		//存进去一个指针，处出来应该是什么？
		//{"get_ptr", mp, args{"set_ptr", &get_int_ptr}, false, get_int_ptr},

		// 存进去一个 nil，用 任意 ptr 取
		// 只有一层的话应该是取不出来了， 毕竟 (*Data)(nil) 接收不了数据
		// 如果是第二层的估计是可以
		//{"get_nil", p, args{"set_nil", (*Data)(nil)}, false, nil},
		{"get_interface", mp, args{"set_interface", &get_int}, false, 1}, //存进去一个 int 的interface , 用 *int 取

		{"get_time", mp, args{"set_time", &get_time}, false, tn},
		{"get_unixTime", mp, args{"set_unixTime", &get_unixTime}, false, UnixTime(tn)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cp.Get(tt.args.key, tt.args.value)

			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryCacheProvider.Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.wantV, reflect.ValueOf(tt.args.value).Elem().Interface()) {
				t.Errorf("MemoryCacheProvider.Get() = %v, want %v", reflect.ValueOf(tt.args.value).Elem().Interface(), tt.wantV)
			}
		})
	}
}

func TestMemoryCacheProvider_MustGet(t *testing.T) {
	MemoryCacheProviderPrepare()
	defer MemoryCacheProviderClearn()

	get_bool := true
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name string
		cp   *MemoryCacheProvider
		args args
	}{
		{"empty_key", mp, args{"", &get_bool}},
		{"empty_ptr", mp, args{"bool_true", (*bool)(nil)}},
		{"not_ptr", mp, args{"bool_false", get_bool}},
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
				tt.cp.MustGet(tt.args.key, tt.args.value)
				return nil
			}()

			if err == nil {
				t.Errorf(tt.name)
			}
		})
	}
}

func TestMemoryCacheProvider_TryGet(t *testing.T) {
	mp.Set("bool_true", true, NoExpiration)

	get_bool := true
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name     string
		cp       *MemoryCacheProvider
		args     args
		wantSucc bool
		wantErr  bool
	}{
		{"get_exists_key", mp, args{"bool_true", &get_bool}, true, false},
		{"get_not_exists_key", mp, args{"bool_true_not", &get_bool}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSucc, err := tt.cp.TryGet(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryCacheProvider.TryGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSucc != tt.wantSucc {
				t.Errorf("MemoryCacheProvider.TryGet() = %v, want %v", gotSucc, tt.wantSucc)
			}
		})
	}
}

func TestMemoryCacheProvider_Create(t *testing.T) {
	mp.Set("bool_true", true, NoExpiration)

	type args struct {
		key   string
		value any
		t     time.Duration
	}
	tests := []struct {
		name    string
		cp      *MemoryCacheProvider
		args    args
		want    bool
		wantErr bool
	}{
		{"create_exists_key", mp, args{"bool_true", true, NoExpiration}, false, false},
		{"create_not_exists_key", mp, args{"bool_true_create", true, 1}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cp.Create(tt.args.key, tt.args.value, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryCacheProvider.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemoryCacheProvider.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryCacheProvider_MustCreate(t *testing.T) {
	type args struct {
		key   string
		value any
		t     time.Duration
	}
	tests := []struct {
		name string
		cp   *MemoryCacheProvider
		args args
	}{
		{"empty_key", mp, args{"", "empty_key", NoExpiration}},
		{"negative_expireTime", mp, args{"negative_expireTime", "negative_expireTime", -1}},
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

				tt.cp.MustCreate(tt.args.key, tt.args.value, tt.args.t)
				return nil
			}()

			if err == nil {
				t.Errorf(tt.name)
			}
		})
	}
}

func TestMemoryCacheProvider_Set(t *testing.T) {
	defer MemoryCacheProviderClearn()

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
		cp      *MemoryCacheProvider
		args    args
		wantErr bool
	}{
		{"set_bool_true", mp, args{"bool_true", true, NoExpiration}, false},
		{"set_bool_false", mp, args{"bool_false", false, NoExpiration}, false},

		{"set_int8", mp, args{"set_int8", int8(-8), NoExpiration}, false},
		{"set_int16", mp, args{"set_int16", int16(-16), NoExpiration}, false},
		{"set_int32", mp, args{"set_int32", int16(-32), NoExpiration}, false},
		{"set_int64", mp, args{"set_int64", int64(-64), NoExpiration}, false},
		{"set_int", mp, args{"set_int", int(100), NoExpiration}, false},

		{"set_uint8", mp, args{"set_uint8", uint8(8), NoExpiration}, false},
		{"set_uint16", mp, args{"set_uint16", uint16(16), NoExpiration}, false},
		{"set_uint32", mp, args{"set_uint32", uint16(32), NoExpiration}, false},
		{"set_uint64", mp, args{"set_uint64", uint64(64), NoExpiration}, false},
		{"set_uint", mp, args{"set_uint", uint(100), NoExpiration}, false},

		{"set_float32", mp, args{"set_float32", float32(123.5), NoExpiration}, false},
		{"set_float64", mp, args{"set_float64", float64(123.5), NoExpiration}, false},

		{"set_array_int", mp, args{"set_array_int", [3]int{1, 2, 3}, NoExpiration}, false},
		{"set_map_int_string", mp, args{"set_map_int_string", map[int]string{1: "a", 2: "b", 3: "c"}, NoExpiration}, false},
		{"set_slice_string", mp, args{"set_slice_string", []string{"a", "b", "c"}, NoExpiration}, false},
		{"set_struct", mp, args{"set_struct", data, NoExpiration}, false},

		{"set_ptr", mp, args{"set_ptr", intvPtr, NoExpiration}, false},
		{"set_nil", mp, args{"set_nil", nil, NoExpiration}, false},
		{"set_interface", mp, args{"set_interface", intf, NoExpiration}, false},

		{"set_time", mp, args{"set_time", tn, NoExpiration}, false},
		{"set_unixTime", mp, args{"set_unixTime", UnixTime(tn), NoExpiration}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cp.Set(tt.args.key, tt.args.value, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("MemoryCacheProvider.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryCacheProvider_MustSet(t *testing.T) {
	type args struct {
		key   string
		value any
		t     time.Duration
	}
	tests := []struct {
		name      string
		cp        *MemoryCacheProvider
		args      args
		wantPanic bool
	}{
		{"empty_key", mp, args{"", "empty_key", NoExpiration}, true},
		{"negative_expireTime", mp, args{"negative_expireTime", "negative_expireTime", -1}, true},
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
				tt.cp.MustSet(tt.args.key, tt.args.value, tt.args.t)
				return nil
			}()

			if tt.wantPanic && err == nil {
				t.Fail()
			}
		})
	}
}

func TestMemoryCacheProvider_Remove(t *testing.T) {
	mp.Set("remov_exists_key", 1, NoExpiration)
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		cp      *MemoryCacheProvider
		args    args
		want    bool
		wantErr bool
	}{
		{"remov_exists_key", mp, args{"remov_exists_key"}, true, false},
		{"remov_not_exists_key", mp, args{"bool_true_Remove"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cp.Remove(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryCacheProvider.Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemoryCacheProvider.Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryCacheProvider_MustRemove(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		cp   *MemoryCacheProvider
		args args
	}{
		{"empty_key", mp, args{""}},
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
				tt.cp.MustRemove(tt.args.key)
				return nil
			}()
			if err == nil {
				t.Errorf(tt.name)
			}
		})
	}
}

func TestMemoryCacheProvider_Increase(t *testing.T) {
	mp.Set("Increase_int64", int64(1), 10*time.Second)      //int64 的数字
	mp.Set("Increase_int8", int8(1), 10*time.Second)        //非 int64 的数字
	mp.Set("Increase_string", "not number", 10*time.Second) //非 int64 的数字

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		cp      *MemoryCacheProvider
		args    args
		want    int64
		wantErr bool
	}{
		{"Increase_int64", mp, args{"Increase_int64"}, 2, false},
		{"Increase_int8", mp, args{"Increase_int8"}, 2, false},
		{"Increase_string", mp, args{"Increase_string"}, 0, true},
		{"Increase_not_exists_key", mp, args{"Increase_not_exists_key"}, 0, true},
		{"Increase_empty_key", mp, args{""}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cp.Increase(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryCacheProvider.Increase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemoryCacheProvider.Increase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryCacheProvider_IncreaseOrCreate(t *testing.T) {
	mp.Set("IncreaseOrCreate_int64", int64(1), 10*time.Second)      //int64 的数字
	mp.Set("IncreaseOrCreate_int8", int8(1), 10*time.Second)        //非 int64 的数字
	mp.Set("IncreaseOrCreate_string", "not number", 10*time.Second) //非 int64 的数字

	type args struct {
		key       string
		increment int64
		t         time.Duration
	}
	tests := []struct {
		name    string
		cp      *MemoryCacheProvider
		args    args
		want    int64
		wantErr bool
	}{
		{"IncreaseOrCreate_int64", mp, args{"IncreaseOrCreate_int64", 2, time.Minute}, 3, false},
		{"IncreaseOrCreate_int8", mp, args{"IncreaseOrCreate_int8", 2, time.Minute}, 3, false},
		{"IncreaseOrCreate_string", mp, args{"IncreaseOrCreate_string", 2, time.Minute}, 0, true},
		{"IncreaseOrCreate_not_exists_key", mp, args{"IncreaseOrCreate_not_exists_key", 2, time.Minute}, 2, false},
		{"IncreaseOrCreate_empty_key", mp, args{"", 2, time.Minute}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cp.IncreaseOrCreate(tt.args.key, tt.args.increment, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryCacheProvider.IncreaseOrCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemoryCacheProvider.IncreaseOrCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryCacheProvider_legalExpireTime(t *testing.T) {
	type args struct {
		t time.Duration
	}
	tests := []struct {
		name string
		m    *MemoryCacheProvider
		args args
		want time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.legalExpireTime(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryCacheProvider.legalExpireTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
