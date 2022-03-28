package cache

import (
	"testing"
	"time"
)

func TestNewExpiration(t *testing.T) {
	type args struct {
		baseExpireTime  time.Duration
		randomRangeTime time.Duration
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{"ok", args{1 * time.Minute, 30 * time.Second}, false},
		{"negative1", args{-1 * time.Minute, 30 * time.Second}, true},
		{"negative2", args{1 * time.Minute, -30 * time.Second}, true},
		{"less_than", args{30 * time.Second, 1 * time.Minute}, true},
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
				NewExpiration(tt.args.baseExpireTime, tt.args.randomRangeTime)
				return nil
			}()
			if (err != nil) != tt.wantPanic {
				t.Errorf(tt.name)
			}
		})
	}
}

func TestNewExpireTimeFrom(t *testing.T) {
	type From func(int64, int64) *Expiration

	type args struct {
		from     From
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		{"Millisecond", args{NewExpirationFromMillisecond, time.Millisecond}},
		{"Second", args{NewExpirationFromSecond, time.Second}},
		{"Minute", args{NewExpirationFromMinute, time.Minute}},
		{"Hour", args{NewExpirationFromHour, time.Hour}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.from(10, 2)
			rangS := 8 * tt.args.duration
			rangE := 12 * tt.args.duration
			for i := 0; i < 1000; i++ {
				next := got.NextExpireTime()
				if next < rangS || next > rangE {
					t.Errorf("get next expire time error[%d - %d]: %d", rangS, rangE, next)
				}
			}
		})
	}

	t.Run("not-rand", func(t *testing.T) {
		var base time.Duration = 10
		exp := NewExpiration(base, 0)
		next := exp.NextExpireTime()

		if base != next {
			t.Fatal("no rand, should be equal")
		}
	})
}
