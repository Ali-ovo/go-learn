package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MyLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()
		ctx.Set("example", "12345")

		// 继续处理请求
		ctx.Next()

		end := time.Since(t)

		fmt.Printf("耗时: %v \r\n", end)
		status := ctx.Writer.Status()

		fmt.Println("status: \r\n", status)
	}
}

func TokenRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string

		for k, v := range ctx.Request.Header {

			if k == "X-Token" {
				token = v[0]
			}

		}

		if token != "123" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg": "未登录",
			})

			ctx.Abort()
		}

		ctx.Next()
	}
}

func main() {

	r := gin.Default()

	// 全局 use
	// r.Use(TokenRequired())

	r.GET("/ping", func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// r.LoadHTMLFiles("template/hello.tmpl")
	r.LoadHTMLGlob("template/*")
	r.GET("/html", func(ctx *gin.Context) {

		ctx.HTML(http.StatusOK, "hello.tmpl", gin.H{
			"title": "hello",
			"msg":   "hello world",
		})
	})

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
