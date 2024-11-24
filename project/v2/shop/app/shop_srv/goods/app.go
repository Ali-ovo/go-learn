package srv

//
//import (
//	"fmt"
//	"github.com/apache/rocketmq-client-go/v2"
//	"github.com/apache/rocketmq-client-go/v2/consumer"
//	"github.com/apache/rocketmq-client-go/v2/primitive"
//	"os"
//	"os/signal"
//	"syscall"
//)
//
//func NewApp(basename string) {
//	//cfg := config.NewConfig()
//	//
//	//return app.
//
//	c, _ := rocketmq.NewPushConsumer(
//		consumer.WithGroupName("goods_canal"),
//		consumer.WithNsResolver(
//			primitive.NewPassthroughResolver([]string{"192.168.101.49:9876"}),
//		),
//	)
//	err := c.Subscribe(
//		"goods_canal",
//		consumer.MessageSelector{},
//		es.GoodsSaveToES,
//	)
//	if err = c.Start(); err != nil {
//		panic(fmt.Sprintf("启动 producer 失败: %s", err))
//	}
//
//	// 接收终止信号
//	quit := make(chan os.Signal)
//	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//	<-quit
//	_ = c.Shutdown()
//}
