package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	routerGroup := r.Group("/goods")
	{
		routerGroup.GET("/list/:name", goodsList)
	}

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func goodsList(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": c.Param("name"),
		"query":   c.DefaultQuery("query", "default"),
		"queryA":  c.DefaultQuery("a", "default"),
	})
}
