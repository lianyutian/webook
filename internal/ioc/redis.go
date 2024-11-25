package ioc

import "github.com/redis/go-redis/v9"

func InitRedis() redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "116.198.217.158:6379",
		Password: "MIIEowIBAAKCAQEAwG90ULRHmAXFXQzZSwleoYts2+bCzUvqhhqtGiv/F5kUsETY", // no password set
		DB:       0,                                                                  // use default DB
	})
	return rdb
}
