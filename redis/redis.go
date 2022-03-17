package redis

import "github.com/go-redis/redis/v8"

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "localhost:6379",
	})

}
