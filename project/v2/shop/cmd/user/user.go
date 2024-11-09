package main

import (
	user "go-learn/project/v2/shop/app/user/srv"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	user.NewApp("user_server").Run()
}
