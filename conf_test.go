package cache

import (
	"github.com/go-redis/redis/v8"
)

var redisC = redis.NewClient(&redis.Options{
	Addr:     "106.52.180.232:6379",
	Password: "1234", // no password set
	DB:       0,      // use default DB
})

var redisCp = NewRedisCacheProvider(redisC)
