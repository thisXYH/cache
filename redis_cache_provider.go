package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis 类型的缓存提供器。
type RedisCacheProvider struct {
	client redis.Cmdable
}

var (
	_ CacheProvider = (*RedisCacheProvider)(nil)
)

func NewRedisCacheProvider(cli redis.Cmdable) *RedisCacheProvider {
	if cli == nil {
		panic(errors.New("param 'cli' is nil"))
	}
	return &RedisCacheProvider{cli}
}

// implement CacheProvider.Get .
func (cli *RedisCacheProvider) Get(key string, value any) error {
	_, err := cli.TryGet(key, value)
	return err
}

// implement CacheProvider.TryGet .
func (cli *RedisCacheProvider) TryGet(key string, value any) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key must not be empty")
	}

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

// implement CacheProvider.Create .
func (cli *RedisCacheProvider) Create(key string, value any, t time.Duration) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key must not be empty")
	}

	v, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	cmd := cli.client.SetNX(context.Background(), key, string(v), t)
	return cmd.Result()
}

// implement CacheProvider.Set .
func (cli *RedisCacheProvider) Set(key string, value any, t time.Duration) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}

	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	const OK = `OK` // 执行成功的返回值。
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

// implement CacheProvider.Remove .
func (cli *RedisCacheProvider) Remove(key string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key must not be empty")
	}

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

// implement CacheProvider.Increase .
func (cli *RedisCacheProvider) Increase(key string) (int64, error) {
	if key == "" {
		return 0, fmt.Errorf("key must not be empty")
	}

	const MaxRetries = 2 // 最大重试次数。

	var value int64 = 0
	increaseIfExistsTrans := func(tx *redis.Tx) error {
		cmd := tx.Get(tx.Context(), key)
		kv, err := cmd.Int64() // 只关心key存不存在，以及是不是数字。
		if err != nil {
			if err == redis.Nil { //缓存不存在。
				return fmt.Errorf("cache key does not exist: %s", key)
			}
			// 存在但不是数字，或者其他 error。
			return err
		}
		kv++
		value = kv
		_, err = tx.TxPipelined(tx.Context(), func(pipe redis.Pipeliner) error {
			// 这边应该默认返回了空的结果，具体执行情况是需要管道执行结束才有的，
			// 所以不需要接收。
			pipe.Incr(tx.Context(), key)
			return nil
		})
		return err
	}

	type Watcher interface {
		Watch(context.Context, func(*redis.Tx) error, ...string) error
	}
	watcher, ok := cli.client.(Watcher)
	if !ok {
		return 0, fmt.Errorf("unsupport redis client type: %t", cli.client)
	}
	for i := 0; i < MaxRetries; i++ {
		err := watcher.Watch(context.Background(), increaseIfExistsTrans, key)
		if err == nil {
			return value, err
		}
		if err == redis.TxFailedErr {
			continue
		}
		return 0, err
	}

	return 0, fmt.Errorf("increment reached maximum number of retries(%d)", MaxRetries)
}

// implement CacheProvider.IncreaseOrCreate .
func (cli *RedisCacheProvider) IncreaseOrCreate(key string, increment int64, t time.Duration) (int64, error) {
	if key == "" {
		return 0, fmt.Errorf("key must not be empty")
	}

	cmd := cli.client.IncrBy(context.Background(), key, increment)
	v, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	// 如果key是新创建的，指定过期时间。
	if v == increment {
		cli.client.Expire(context.Background(), key, t)
	}

	return v, err
}
