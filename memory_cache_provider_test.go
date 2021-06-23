package cache

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestMemoryCacheProvider(t *testing.T) {
	cache := NewMemoryCacheProvider(10)

	cache.Set("ffaf", int64(10), time.Minute)

	var i int8
	cache.MustGet("ffaf", &i)
	fmt.Println(i)

}

func TestNewMemoryCacheProvider(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name string
		args args
		want *MemoryCacheProvider
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemoryCacheProvider(tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemoryCacheProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryCacheProvider_Get(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name    string
		cp      *MemoryCacheProvider
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cp.Get(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("MemoryCacheProvider.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryCacheProvider_MustGet(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name string
		cp   *MemoryCacheProvider
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cp.MustGet(tt.args.key, tt.args.value)
		})
	}
}

func TestMemoryCacheProvider_TryGet(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name    string
		cp      *MemoryCacheProvider
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cp.TryGet(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryCacheProvider.TryGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemoryCacheProvider.TryGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryCacheProvider_Create(t *testing.T) {
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
		// TODO: Add test cases.
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
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cp.MustCreate(tt.args.key, tt.args.value, tt.args.t); got != tt.want {
				t.Errorf("MemoryCacheProvider.MustCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryCacheProvider_Set(t *testing.T) {
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
		// TODO: Add test cases.
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
		name string
		cp   *MemoryCacheProvider
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cp.MustSet(tt.args.key, tt.args.value, tt.args.t)
		})
	}
}

func TestMemoryCacheProvider_Remove(t *testing.T) {
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
		// TODO: Add test cases.
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
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cp.MustRemove(tt.args.key); got != tt.want {
				t.Errorf("MemoryCacheProvider.MustRemove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryCacheProvider_expireIfNeeded(t *testing.T) {
	type args struct {
		key string
		v   memoryCacheData
	}
	tests := []struct {
		name string
		cp   *MemoryCacheProvider
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cp.expireIfNeeded(tt.args.key, tt.args.v); got != tt.want {
				t.Errorf("MemoryCacheProvider.expireIfNeeded() = %v, want %v", got, tt.want)
			}
		})
	}
}
