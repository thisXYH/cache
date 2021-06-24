package cache

import (
	"github.com/go-redis/redis/v8"
)

var redisC = redis.NewClient(&redis.Options{
	Addr:     "host:port",
	Password: "password", // no password set
	DB:       0,          // use default DB
})

var redisCp = NewRedisCacheProvider(redisC)
