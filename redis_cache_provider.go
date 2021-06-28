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
// 数据的组织方式，基础类型直接使用
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

// implement ICacheProvider.Increase
func (cli *RedisCacheProvider) Increase(key string) (int64, error) {
	const MaxRetries = 3 //最大重试次数。

	var value int64 = 0
	increaseIfExistsTrans := func(tx *redis.Tx) error {
		cmd := tx.Get(tx.Context(), key)
		kv, err := cmd.Int64() //只关心key存不存在，以及是不是数字。
		if err != nil {
			if err == redis.Nil { //缓存不存在
				return fmt.Errorf("cache key not exists")
			}
			// 存在但不是数字，或者其他error
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

// implement ICacheProvider.MustIncrease
func (cli *RedisCacheProvider) MustIncrease(key string) int64 {
	v, err := cli.Increase(key)
	if err != nil {
		panic(err)
	}
	return v
}

// implement ICacheProvider.IncreaseOrCreate
func (cli *RedisCacheProvider) IncreaseOrCreate(key string, increment int64, t time.Duration) (int64, error) {
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

// implement ICacheProvider.MustIncreaseOrCreate
func (cli *RedisCacheProvider) MustIncreaseOrCreate(key string, increment int64, t time.Duration) int64 {
	v, err := cli.IncreaseOrCreate(key, increment, t)
	if err != nil {
		panic(err)
	}
	return v
}
