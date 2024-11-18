package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	for {
		ops.Inc()
		time.Sleep(2 * time.Second)
	}
}

var (
	ops = promauto.NewCounter(prometheus.CounterOpts{
		Name: "shop_test",
		Help: "just for test",
	})
)

func main() {
	go recordMetrics()
	r := gin.Default()
	r.GET("/metrics", gin.WrapH(promhttp.Handler())) // promhttp.Handler() 这里已经自动做好了 获取基本信息的逻辑处理
	r.Run("0.0.0.0:8050")
}
