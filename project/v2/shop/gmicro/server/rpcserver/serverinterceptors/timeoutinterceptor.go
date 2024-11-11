package serverinterceptors

import (
	"context"
	"fmt"
	"shop/gmicro/pkg/code"
	"shop/gmicro/pkg/errors"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// UnaryTimeoutInterceptor returns a func that sets timeout to incoming unary requests.  返回一个函数 为一元请求 设置超时时间
func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var resp any // 2个 goroutine 操作全局变量 需要加锁  因为会并发地读写同一个全局变量
		var err error
		var lock sync.Mutex
		done := make(chan struct{})
		// create channel with buffer size 1 to avoid goroutine leak
		panicChan := make(chan any, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					// attach call stack to avoid missing in different goroutine
					panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack())))
				}
			}()

			lock.Lock()
			defer lock.Unlock()
			resp, err = handler(ctx, req)
			close(done) // 关闭通道
		}()

		select {
		case p := <-panicChan: // 错误 抛出
			panic(p)
		case <-done: // done 关闭 说明运行成功  执行以下逻辑
			lock.Lock()
			defer lock.Unlock()
			return resp, err
		case <-ctx.Done(): // 超时了 弹出错误信息
			err := ctx.Err()
			if err == context.Canceled {
				err = errors.WithCode(code.ErrCanceledGrpc, err.Error())
				//err = status.Error(codes.Canceled, err.Error())
			} else if err == context.DeadlineExceeded {
				err = errors.WithCode(code.ErrDeadlineExceededGrpc, err.Error())
				//err = status.Error(codes.DeadlineExceeded, err.Error())

			}
			return nil, errors.ToGrpcError(err)
		}
	}
}
