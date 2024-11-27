package main

import (
	"context"
	"fmt"
	inventory_pb "shop/api/inventory/v1"
	"shop/gmicro/registry/consul"
	rpc "shop/gmicro/server/rpcserver"
	"shop/gmicro/server/rpcserver/clientinterceptors"
	"shop/gmicro/server/rpcserver/selector"
	"shop/gmicro/server/rpcserver/selector/random"
	"time"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
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
					Num:     11,
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
			return
		}
		c.JSON(200, gin.H{"message": "ok"})
		return
	})
	r.GET("/inventory", func(c *gin.Context) {
		selector.SetGlobalSelector(random.NewBuilder())
		rpc.InitBuilder()

		conf := api.DefaultConfig()
		conf.Address = "172.16.89.133:8500"
		conf.Scheme = "http"
		cli, err := api.NewClient(conf)
		if err != nil {
			panic(err)
		}

		conn, err := rpc.DialInsecure(
			context.Background(),
			rpc.WithBanlancerName("selector"),
			rpc.WithDiscovery(consul.New(cli)),
			rpc.WithEndpoint("discovery:///shop-inventory-srv"),
			rpc.WithClientEnableTracing(false),
			rpc.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
			rpc.WithClientTimeout(time.Second*5000),
		)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		uc := inventory_pb.NewInventoryClient(conn)

		re, err := uc.InvDetail(context.Background(), &inventory_pb.GoodsInvInfo{
			GoodsId: 421,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(re)
		return
	})
	r.Run(":8089")
}
