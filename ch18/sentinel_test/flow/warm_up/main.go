package main

import (
	"fmt"
	"log"
	"math/rand/v2"
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
			TokenCalculateStrategy: flow.WarmUp, // 冷启动策略
			ControlBehavior:        flow.Reject, // 直接拒绝
			Threshold:              1000,
			WarmUpPeriodSec:        60, // 预热时间
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}
	ch := make(chan struct{})

	var globalTotal int
	var passTotal int
	var blockTotal int
	for i := 0; i < 100; i++ {
		go func() {
			for {
				globalTotal++
				e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					blockTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					passTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					e.Exit()
				}
			}
		}()
	}

	go func() {
		var oldTotal int
		var oldPass int
		var oldBlock int
		for {
			oneSecondTotal := globalTotal - oldTotal
			oldTotal = globalTotal

			oneSecondPass := passTotal - oldPass
			oldPass = passTotal

			oneSecondBlock := blockTotal - oldBlock
			oldBlock = blockTotal

			time.Sleep(time.Second)
			// fmt.Printf("globalTotal: %d, passTotal: %d, blockTotal: %d\n", globalTotal, passTotal, blockTotal)
			fmt.Printf("oneSecondTotal: %d, oneSecondPass: %d, oneSecondBlock: %d\n", oneSecondTotal, oneSecondPass, oneSecondBlock)
		}
	}()

	<-ch
}
