package main

import (
	"sync"
	"time"
)

func main() {

	var rwlock sync.RWMutex
	var wg sync.WaitGroup

	wg.Add(6)

	// write
	go func() {
		time.Sleep(time.Second)
		defer wg.Done()
		rwlock.Lock() // 写锁 会防止其他写锁和读锁
		defer rwlock.Unlock()
		println("get write lock")
		time.Sleep(time.Second * 5)
	}()

	for i := 0; i < 5; i++ {
		// read
		go func() {
			defer wg.Done()

			for {
				rwlock.RLock()
				time.Sleep(500 * time.Millisecond)
				println("get read lock")
				rwlock.RUnlock()
			}
		}()
	}

	wg.Wait()
}
