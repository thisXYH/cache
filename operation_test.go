package cache

import (
	"testing"
	"time"
)

func TestOperation_oneKeyToStr(t *testing.T) {
	op := &Operation{}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name      string
		c         *Operation
		args      args
		want      string
		wantPanic bool
	}{

		{"bool_false", op, args{false}, "0", false},
		{"bool_true", op, args{true}, "1", false},

		{"int8", op, args{int8(2)}, "2", false},
		{"int16", op, args{int16(3)}, "3", false},
		{"int32", op, args{int32(4)}, "4", false},
		{"int64", op, args{int64(5)}, "5", false},
		{"int", op, args{int(6)}, "6", false},

		{"uint8", op, args{uint8(7)}, "7", false},
		{"uint16", op, args{uint16(8)}, "8", false},
		{"uint32", op, args{uint32(9)}, "9", false},
		{"uint64", op, args{uint64(10)}, "10", false},
		{"uint", op, args{uint(11)}, "11", false},

		{"float32", op, args{float32(12.12)}, "12.12", false},
		{"float64", op, args{float32(13.13)}, "13.13", false},

		{"complex", op, args{complex(14, 14)}, "(14+14i)", false},

		{"string", op, args{"string"}, "string", false},
		{"time", op, args{tn}, UnixTime(tn).String(), false},
		{"unixTime", op, args{UnixTime(tn)}, UnixTime(tn).String(), false},

		// 支持类型的 pointer
		{"pointer", op, args{&tn}, UnixTime(tn).String(), false},

		// panic
		{"nil", op, args{nil}, "", true},
		{"map", op, args{map[string]int{"a": 1}}, "", true},
		{"slice", op, args{[]int{1, 2, 3}}, "", true},
		{"array", op, args{[...]int{1, 2, 3}}, "", true},
		{"channel", op, args{make(chan int)}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := func() (g string, e error) {
				defer func() {
					err := recover()
					if err != nil {
						e = err.(error)
					}
				}()
				g = oneKeyToStr(tt.args.v)
				return
			}()

			if (err != nil) != tt.wantPanic {
				t.Errorf("Operation error = %v, wantP %v", err, tt.wantPanic)
			}

			if err != nil {
				return
			}

			if got != tt.want {
				t.Errorf("Operation.oneKeyToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperation0(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)

	const ns = "ns"
	const prefix = "prefix"

	t.Run("string", func(t *testing.T) {
		op := NewOperation0[string](ns, prefix, provider, CacheExpirationZero)
		k := op.Key()

		var res string
		res, ok, _ := k.TryGet()
		if ok {
			t.Fatal("key should not be")
		}

		k.Set("vv")
		res, ok, _ = k.TryGet()
		if !ok {
			t.Fatal("key should   be")
		}

		if res != "vv" {
			t.Fatal("value mismatch")
		}
	})
}

func TestOperation1(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)

	const ns = "ns"
	const prefix = "prefix"

	t.Run("int-string", func(t *testing.T) {
		op := NewOperation1[int, string](ns, prefix, provider, CacheExpirationZero)
		k := op.Key(100)

		res, ok, _ := k.TryGet()
		if ok {
			t.Fatal("key should not be")
		}

		k.Set("vv")
		res, ok, _ = k.TryGet()
		if !ok {
			t.Fatal("key should   be")
		}

		if res != "vv" {
			t.Fatal("value mismatch")
		}
	})

	t.Run("string-int", func(t *testing.T) {
		op := NewOperation1[string, int](ns, prefix, provider, CacheExpirationZero)
		k := op.Key("g")

		res, ok, _ := k.TryGet()
		if ok {
			t.Fatal("key should not be")
		}

		k.Set(33)
		res, ok, _ = k.TryGet()
		if !ok {
			t.Fatal("key should   be")
		}

		if res != 33 {
			t.Fatal("value mismatch")
		}
	})
}

func TestOperation2(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)

	const ns = "ns"
	const prefix = "prefix"

	t.Run("int-string-float64", func(t *testing.T) {
		op := NewOperation2[int, string, float64](ns, prefix, provider, CacheExpirationZero)
		k := op.Key(100, "gg")

		res, ok, _ := k.TryGet()
		if ok {
			t.Fatal("key should not be")
		}

		k.Set(0.5)
		res, ok, _ = k.TryGet()
		if !ok {
			t.Fatal("key should   be")
		}

		if res != 0.5 {
			t.Fatal("value mismatch")
		}
	})
}

func TestOperation3(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)

	const ns = "ns"
	const prefix = "prefix"

	t.Run("int-string-float64-time", func(t *testing.T) {
		op := NewOperation3[int, string, float64, time.Time](ns, prefix, provider, CacheExpirationZero)
		k := op.Key(100, "gg", 1.5)

		res, ok, _ := k.TryGet()
		if ok {
			t.Fatal("key should not be")
		}

		v := time.Date(2022, 3, 15, 12, 22, 30, 0, time.UTC)
		k.Set(v)
		res, ok, _ = k.TryGet()
		if !ok {
			t.Fatal("key should   be")
		}

		if res != v {
			t.Fatal("value mismatch")
		}
	})
}

func TestOperation4(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)

	const ns = "ns"
	const prefix = "prefix"

	t.Run("int-int-string-float64-string", func(t *testing.T) {
		op := NewOperation4[int, int, string, float64, string](ns, prefix, provider, CacheExpirationZero)
		k := op.Key(100, 200, "gg", 0.5)

		var res string
		res, ok, _ := k.TryGet()
		if ok {
			t.Fatal("key should not be")
		}

		k.Set("vv")
		res, ok, _ = k.TryGet()
		if !ok {
			t.Fatal("key should   be")
		}

		if res != "vv" {
			t.Fatal("value mismatch")
		}
	})
}
