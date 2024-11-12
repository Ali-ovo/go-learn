package main

import (
	"context"
	"fmt"
	v1 "shop/api/user/v1"
	rpc "shop/gmicro/server/rpcserver"
	_ "shop/gmicro/server/rpcserver/resolver/direct" // 必填
)

func main() {
	conn, err := rpc.DialInsecure(context.Background(),
		// 多添加一个 /  因为 方便做切割 direct:///192.168.16.154:8081 转换成 URL.Path: /192.168.16.154:8081  URL.Scheme: direct
		// rpc.WithEndpoint("127.0.0.1:8081"),
		rpc.WithEndpoint("172.16.100.208:8081"),
		// rpc.WithClientTimeout(1*time.Second),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	uc := v1.NewUserClient(conn)
	r, err := uc.GetUserList(context.Background(), &v1.PageInfo{})
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}
