package main

import (
	"context"
	"fmt"
	"go-learn/ch14/stream_grpc/proto"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	c := proto.NewGreeterClient(conn)

	// 服务端流模式
	// res, _ := c.GetStream(context.Background(), &proto.StreamReqData{
	// 	Data: "ali",
	// })

	// for {
	// 	r, err := res.Recv()
	// 	if err != nil {
	// 		break
	// 	}
	// 	fmt.Println(r.GetData())
	// }

	// 客户端流模式
	// putS, _ := c.PutStream(context.Background())
	// i := 0
	// for ; i < 10; i++ {
	// 	putS.Send(&proto.StreamReqData{
	// 		Data: fmt.Sprintf("%d", i),
	// 	})

	// 	time.Sleep(time.Second)
	// }
	// putS.CloseSend()

	// 双向流
	allStr, _ := c.AllStream(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			data, _ := allStr.Recv()
			fmt.Println("客户端 收到数据", data)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			allStr.Send(&proto.StreamReqData{
				Data: fmt.Sprintf("客户端  发送数据 %v %d", time.Now().Unix(), i),
			})
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()

}
