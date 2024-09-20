package main

import (
	"fmt"
	"go-learn/ch14/stream_grpc/proto"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

const PORT = ":8080"

type server struct {
	proto.UnimplementedGreeterServer
}

func (s *server) GetStream(req *proto.StreamReqData, res proto.Greeter_GetStreamServer) error {

	i := 0
	for {
		i++
		res.Send(&proto.StreamResData{
			Data: fmt.Sprintf("%v %d", time.Now().Unix(), i),
		})

		time.Sleep(time.Second)
		if i > 10 {
			break
		}
	}

	return nil
}

func (s *server) PutStream(cliStr proto.Greeter_PutStreamServer) error {

	for {
		if a, err := cliStr.Recv(); err != nil {
			fmt.Println(err)

			break
		} else {
			fmt.Println(a)

		}
	}

	return nil
}

func (s *server) AllStream(allStr proto.Greeter_AllStreamServer) error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			data, _ := allStr.Recv()
			fmt.Println("服务端 收到数据", data)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			allStr.Send(&proto.StreamResData{
				Data: fmt.Sprintf("服务端  发送数据 %v %d", time.Now().Unix(), i),
			})
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()

	return nil
}

func main() {
	listen, err := net.Listen("tcp", PORT)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &server{})
	err = s.Serve(listen)
	if err != nil {
		panic(err)
	}
}
