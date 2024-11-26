package main

import (
	"math/rand"
	"os"
	"runtime"
	"shop/app/shop/api"
	"time"
)

func main() {
	// -c ./configs/shop/api.yaml
	// --log.level=error

	rand.NewSource(time.Now().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 { // 设置的环境变量 GOMAXPROCS 为 0 即设置 本机cpu核数
		runtime.GOMAXPROCS(runtime.NumCPU()) // 用于指定并发执行时可以使用的最大 CPU 核心数。这个函数通常在初始化程序时使用，以确保最大化利用系统资源
	}
	api.NewApp("shop").Run()
}
