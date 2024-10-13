package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

func main() {
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.189.128:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.
	gNum := 2
	mutexname := "421"

	var wg sync.WaitGroup
	wg.Add(gNum)
	for i := 0; i < gNum; i++ {
		go func() {
			defer wg.Done()
			mutex := rs.NewMutex(mutexname)

			fmt.Println("try to get lock")

			if err := mutex.Lock(); err != nil {
				panic(err)
			}

			fmt.Println("get lock success")

			time.Sleep(time.Second * 5)

			fmt.Println("start to unlock")
			if ok, err := mutex.Unlock(); !ok || err != nil {
				panic("unlock failed")
			}
			fmt.Println("unlock success")
		}()
	}

	wg.Wait()
}
