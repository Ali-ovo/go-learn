package main

import (
	"fmt"
	inventory_pb "shop/api/inventory/v1"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
)

func main() {
	r := gin.Default()
	r.GET("/start", func(c *gin.Context) {
		orderSn := shortuuid.New()
		req := &inventory_pb.SellInfo{
			GoodsInfo: []*inventory_pb.GoodsInvInfo{
				{
					GoodsId: 421,
					Num:     2,
				},
			},
			OrderSn: orderSn,
		}
		dtmServer := "172.16.89.133:36790"
		qsBusi := "discovery://172.16.89.133:8500/inventory_srv"
		fmt.Println(orderSn)
		saga := dtmgrpc.NewSagaGrpc(dtmServer, orderSn).
			// 添加一个TransOut的子事务，正向操作为url: qsBusi+"/TransOut"， 逆向操作为url: qsBusi+"/TransOutCom"
			Add(qsBusi+"/Inventory/Sell", qsBusi+"/Inventory/Reback", req)
		// 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
		err := saga.Submit()
		if err != nil {
			fmt.Printf("saga 提交失败: %s\r\n", err.Error())
			c.JSON(500, gin.H{"message": err.Error()})
		}
		c.JSON(200, gin.H{"message": "ok"})
	})
	r.Run(":8089")
}
