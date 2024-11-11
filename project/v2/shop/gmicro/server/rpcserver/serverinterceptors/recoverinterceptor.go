package serverinterceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop/gmicro/pkg/log"
	"runtime/debug"
)

// PS: 参考 go-zero

// StreamRecoverInterceptor catches panics in processing stream requests and recovers.
func StreamRecoverInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r any) {
		err = toPanicError(r)
	})
	return handler(srv, ss)
}

// UnaryRecoverInterceptor catches panics in processing unary requests and recovers.
func UnaryRecoverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer handleCrash(func(r any) {
		err = toPanicError(r)
	})
	return handler(ctx, req)
}

func handleCrash(handler func(any)) {
	if r := recover(); r != nil {
		handler(r)
	}
}

func toPanicError(r any) error {
	log.Errorf("%+v\n\n%s", r, debug.Stack())
	return status.Errorf(codes.Internal, "panic: %v", r)
}
