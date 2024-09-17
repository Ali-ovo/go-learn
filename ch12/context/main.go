package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

// var stop = make(chan struct{})

// func cpuInfo(stop chan struct{}) {
func cpuInfo(ctx context.Context) {
	// 获取 id
	fmt.Printf("tracId:%s\r\n", ctx.Value("traceId"))

	defer wg.Done()
	for {

		select {
		case <-ctx.Done():
			fmt.Println("stop")
			return

		default:
			time.Sleep(2 * time.Second)
			fmt.Println("cpu")
		}
	}
}

func main() {

	// stop := make(chan struct{})

	wg.Add(1)

	// WithCancel 手动取消
	// ctx, cancel := context.WithCancel(context.Background())
	// ctx1, _ := context.WithCancel(ctx)

	// WithTimeout 自动超时
	ctx, _ := context.WithTimeout(context.Background(), 6*time.Second)

	// WithDeadline 在指定的时间 cancel,也是 WithTimeout 的原理

	// WithValue 增加值
	valueCtx := context.WithValue(ctx, "traceId", "test123")

	// go cpuInfo(stop)
	go cpuInfo(valueCtx)

	// time.Sleep(5 * time.Second)

	// stop <- struct{}{}
	// cancel() // WithTimeout auto call cancel

	wg.Wait()

	fmt.Println("监控完成")

}
