package caching

import (
	"testing"
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
				g = tt.c.oneKeyToStr(tt.args.v)
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
