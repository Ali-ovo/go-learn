package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", "192.168.189.128", 6379),
	})

	value, err := rdb.Get(context.Background(), "18375713787").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(value)
}
