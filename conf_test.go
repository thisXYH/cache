package cache_test

import "github.com/go-redis/redis/v8"

var RedisClient = redis.NewClient(&redis.Options{
	Addr:     "HOST:PORT",
	Password: "PASSWORD", // no password set
	DB:       0,          // use default DB
})
