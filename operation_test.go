package cache

import (
	"testing"
	"time"
)

func Test_oneKeyToStr(t *testing.T) {
	type args struct {
		v interface{}
	}

	tests := []struct {
		name      string
		args      args
		want      string
		wantPanic bool
	}{
		{"bool_false", args{false}, "0", false},
		{"bool_true", args{true}, "1", false},

		{"int8", args{int8(2)}, "2", false},
		{"int16", args{int16(3)}, "3", false},
		{"int32", args{int32(4)}, "4", false},
		{"int64", args{int64(5)}, "5", false},
		{"int", args{int(6)}, "6", false},

		{"uint8", args{uint8(7)}, "7", false},
		{"uint16", args{uint16(8)}, "8", false},
		{"uint32", args{uint32(9)}, "9", false},
		{"uint64", args{uint64(10)}, "10", false},
		{"uint", args{uint(11)}, "11", false},

		{"float32", args{float32(12.12)}, "12.12", false},
		{"float64", args{float32(13.13)}, "13.13", false},

		{"complex", args{complex(14, 14)}, "(14+14i)", false},

		{"string", args{"string"}, "string", false},
		{"time", args{tn}, UnixTime(tn).String(), false},
		{"unixTime", args{UnixTime(tn)}, UnixTime(tn).String(), false},

		// 支持类型的 pointer
		{"pointer", args{&tn}, UnixTime(tn).String(), false},

		// panic
		{"nil", args{nil}, "", true},
		{"map", args{map[string]int{"a": 1}}, "", true},
		{"slice", args{[]int{1, 2, 3}}, "", true},
		{"array", args{[...]int{1, 2, 3}}, "", true},
		{"channel", args{make(chan int)}, "", true},
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

		if op.Key().Key != "ns:prefix" {
			t.Fatal("key error")
		}
	})
}

func TestOperation1(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)

	const ns = "ns"
	const prefix = "prefix"

	t.Run("bool-string", func(t *testing.T) {
		op := NewOperation1[bool, string](ns, prefix, provider, CacheExpirationZero)

		if op.Key(true).Key != "ns:prefix_1" {
			t.Fatal("key error")
		}
	})
}

func TestOperation2(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)

	const ns = "ns"
	const prefix = "prefix"

	t.Run("bool-int-string", func(t *testing.T) {
		op := NewOperation2[bool, int, string](ns, prefix, provider, CacheExpirationZero)

		if op.Key(true, 9).Key != "ns:prefix_1_9" {
			t.Fatal("key error")
		}
	})
}

func TestOperation3(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)

	const ns = "ns"
	const prefix = "prefix"

	t.Run("bool-int-uint-string", func(t *testing.T) {
		op := NewOperation3[bool, int, uint, string](ns, prefix, provider, CacheExpirationZero)

		if op.Key(true, 9, 8).Key != "ns:prefix_1_9_8" {
			t.Fatal("key error")
		}
	})
}

func TestOperation4(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)

	const ns = "ns"
	const prefix = "prefix"

	t.Run("bool-int-uint-float32-string", func(t *testing.T) {
		op := NewOperation4[bool, int, uint, float32, string](ns, prefix, provider, CacheExpirationZero)

		if op.Key(true, 9, 8, 0.50).Key != "ns:prefix_1_9_8_0.5" {
			t.Fatal("key error")
		}
	})
}

func TestOperation5(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)

	const ns = "ns"
	const prefix = "prefix"

	t.Run("bool-int-uint-float32-byte-string", func(t *testing.T) {
		op := NewOperation5[bool, int, uint, float32, byte, string](ns, prefix, provider, CacheExpirationZero)

		if op.Key(true, 9, 8, 0.50, 255).Key != "ns:prefix_1_9_8_0.5_255" {
			t.Fatal("key error")
		}
	})
}

func TestOperation6(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)

	const ns = "ns"
	const prefix = "prefix"

	t.Run("bool-int-uint-float32-byte-string-string", func(t *testing.T) {
		op := NewOperation6[bool, int, uint, float32, byte, string, string](ns, prefix, provider, CacheExpirationZero)

		if op.Key(true, 9, 8, 0.50, 255, "a").Key != "ns:prefix_1_9_8_0.5_255_a" {
			t.Fatal("key error")
		}
	})
}

func TestOperation7(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)

	const ns = "ns"
	const prefix = "prefix"

	t.Run("bool-int-uint-float32-byte-string-time-string", func(t *testing.T) {
		op := NewOperation7[bool, int, uint, float32, byte, string, time.Time, string](ns, prefix, provider, CacheExpirationZero)

		tv := time.Date(2022, 03, 29, 15, 16, 00, 00, time.UTC)
		unixTv := UnixTime(tv)

		if op.Key(true, 9, 8, 0.50, 255, "a", tv).Key != "ns:prefix_1_9_8_0.5_255_a_"+unixTv.String() {
			t.Fatal("key error")
		}
	})
}

func TestOperation8(t *testing.T) {
	provider := NewMemoryCacheProvider(time.Second)

	const ns = "ns"
	const prefix = "prefix"

	t.Run("bool-int-uint-float32-byte-string-time-unixTime-string", func(t *testing.T) {
		op := NewOperation8[bool, int, uint, float32, byte, string, time.Time, UnixTime, string](ns, prefix, provider, CacheExpirationZero)

		tv := time.Date(2022, 03, 29, 15, 16, 00, 00, time.UTC)
		unixTv := UnixTime(tv)

		if op.Key(true, 9, 8, 0.50, 255, "a", tv, unixTv).Key != "ns:prefix_1_9_8_0.5_255_a_"+unixTv.String()+"_"+unixTv.String() {
			t.Fatal("key error")
		}
	})
}
