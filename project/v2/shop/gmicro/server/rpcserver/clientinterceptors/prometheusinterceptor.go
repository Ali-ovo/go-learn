package clientinterceptors

import (
	"context"
	"shop/gmicro/core/metric"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const serverNamespace = "rpc_client"

/*
两个基本指标: 1. 每个请求的耗时( histogram ) 2. 每个请求的状态计数器( counter )
/user 状态码 有label 主要是状态码
*/

var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "mxshop_duration_ms",
		Help:      "rpc server requests duration(ms)",
		Labels:    []string{"method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "shop_code_total",
		Help:      "rpc server requests code count.",
		Labels:    []string{"method", "code"},
	})
)

func UnaryPrometheusInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	startTime := time.Now()
	// 记录了耗时
	err := invoker(ctx, method, req, reply, cc, opts...)
	// 参数 截至时间 以 毫秒为单位 rpc 的方法名作为标签
	metricServerReqDur.Observe(int64(time.Since(startTime)/time.Millisecond), method)
	// 记录状态码 和 访问次数
	metricServerReqCodeTotal.Inc(method, strconv.Itoa(int(status.Code(err))))
	return err
}
