package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis 类型的缓存提供器。
// 数据的组织方式，基础类型直接使用
type RedisCacheProvider struct {
	client *redis.Client
}

var (
	_ ICacheProvider = (*RedisCacheProvider)(nil)
)

func NewRedisCacheProvider(cli *redis.Client) *RedisCacheProvider {
	if cli == nil {
		panic("param 'cli' is nil")
	}
	return &RedisCacheProvider{cli}
}

// implement ICacheProvider.Get
func (cli *RedisCacheProvider) Get(key string, value any) error {
	_, err := cli.TryGet(key, value)
	return err
}

// implement ICacheProvider.MustGet
func (cli *RedisCacheProvider) MustGet(key string, value any) {
	if err := cli.Get(key, value); err != nil {
		panic(err)
	}
}

// implement ICacheProvider.TryGet
func (cli *RedisCacheProvider) TryGet(key string, value any) (bool, error) {
	cmd := cli.client.Get(context.Background(), key)
	v, err := cmd.Result()
	if err != nil {
		if err == redis.Nil { //key 不存在
			return false, nil
		}
		return false, err
	}

	if err = json.Unmarshal([]byte(v), value); err != nil {
		return false, err
	}

	return true, nil
}

// implement ICacheProvider.Create
func (cli *RedisCacheProvider) Create(key string, value any, t time.Duration) (bool, error) {
	v, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	cmd := cli.client.SetNX(context.Background(), key, string(v), t)
	return cmd.Result()
}

// implement ICacheProvider.MustCreate
func (cli *RedisCacheProvider) MustCreate(key string, value any, t time.Duration) bool {
	v, err := cli.Create(key, value, t)
	if err != nil {
		panic(err)
	}
	return v
}

// implement ICacheProvider.Set
func (cli *RedisCacheProvider) Set(key string, value any, t time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	const OK = `OK` //执行成功的返回值。
	cmd := cli.client.Set(context.Background(), key, string(v), t)
	cv, err := cmd.Result()
	if err != nil {
		return err
	}

	if cv != OK {
		return errors.New("set key failed:" + cv)
	}

	return nil
}

// implement ICacheProvider.MustSet
func (cli *RedisCacheProvider) MustSet(key string, value any, t time.Duration) {
	if err := cli.Set(key, value, t); err != nil {
		panic(err)
	}
}

// implement ICacheProvider.Remove
func (cli *RedisCacheProvider) Remove(key string) (bool, error) {
	cmd := cli.client.Del(context.Background(), key)
	v, err := cmd.Result()
	if err != nil {
		panic(err)
	}
	if v <= 0 {
		return false, nil
	}

	return true, nil
}

// implement ICacheProvider.MustRemove
func (cli *RedisCacheProvider) MustRemove(key string) bool {
	v, err := cli.Remove(key)
	if err != nil {
		panic(err)
	}
	return v
}
