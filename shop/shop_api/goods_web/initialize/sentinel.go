package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"go.uber.org/zap"
)

func InitSentinel() {
	conf := config.NewDefaultConfig()
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()

	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Fatalf("Unexpected error: %+v", err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "goods-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			// 6 秒内最多 3 个请求
			Threshold:        3,
			StatIntervalInMs: 6000,
		},
	})
	if err != nil {
		zap.S().Fatalf("Unexpected error: %+v", err)
		return
	}
}
