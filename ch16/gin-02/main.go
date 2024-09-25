package main

import (
	"net/http"

	"go-learn/ch16/gin-02/proto"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/json", moreJSON)
	r.GET("/someProtoBuf", returnProto)

	r.Run("0.0.0.0:8081") // 监听并在 0.0.0.0:8080 上启动服务
}

func moreJSON(c *gin.Context) {
	var msg struct {
		Name    string `json:"user"`
		Message string
		Number  int
	}

	msg.Name = "root"
	msg.Message = "test json"
	msg.Number = 200

	c.JSON(http.StatusOK, msg)

}

func returnProto(c *gin.Context) {
	course := []string{"go", "java", "web"}
	user := &proto.Teacher{
		Name:   "root",
		Course: course,
	}

	c.ProtoBuf(http.StatusOK, user)
}
