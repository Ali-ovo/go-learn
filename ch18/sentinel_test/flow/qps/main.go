package main

import (
	"fmt"
	"log"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
)

func main() {
	conf := config.NewDefaultConfig()
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()

	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling, // 迅速通过
			Threshold:              100,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "some-test2",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject, // 直接拒绝
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}

	for i := 0; i < 12; i++ {
		e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			// time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
			fmt.Println("Blocked")
		} else {
			fmt.Println("Passed")
			e.Exit()
		}
		time.Sleep(11 * time.Millisecond)

	}
}
