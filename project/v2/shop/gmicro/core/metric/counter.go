package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type (
	// A CounterVecOpts is an alias of VectorOpts.
	CounterVecOpts VectorOpts

	// CounterVec interface represents a counter vector.
	CounterVec interface {
		// Inc increments labels.
		Inc(labels ...string)
		// Add adds labels with v.
		Add(v float64, labels ...string)
	}

	promCounterVec struct {
		// CounterVec允许在同一个指标名称下记录多个不同维度的计数器值，每个维度都有自己的标签，更加灵活
		// 用于记录只增不减的度量数据，即Counter类型指标。
		// Counter指标通常用于表示计数器的数量，例如HTTP请求数、错误次数、任务执行次数等。
		// CounterVec同样可以包含多个维度，每个维度对应一个标签，方便对计数器进行细粒度的区分和查询。
		counter *prom.CounterVec
	}
)

// NewCounterVec returns a CounterVec.
func NewCounterVec(cfg *CounterVecOpts) CounterVec {
	if cfg == nil {
		return nil
	}

	// 这里有多种声明方式 也可以这样声明
	//vec := promauto.NewCounterVec(prom.CounterOpts{
	//	Namespace: cfg.Namespace, // 指标的命名空间 并避免不同应用程序之间的指标名称冲突
	//	Subsystem: cfg.Subsystem, // 指标的子系统名称 用于更进一步组织指标，并区分同一应用程序中的不同组件
	//	Name:      cfg.Name,      // 指标的名称 唯一
	//	Help:      cfg.Help,      // 指标的相关信息
	//}, cfg.Labels)
	vec := prom.NewCounterVec(prom.CounterOpts{
		Namespace: cfg.Namespace, // 指标的命名空间 并避免不同应用程序之间的指标名称冲突
		Subsystem: cfg.Subsystem, // 指标的子系统名称 用于更进一步组织指标，并区分同一应用程序中的不同组件
		Name:      cfg.Name,      // 指标的名称 唯一
		Help:      cfg.Help,      // 指标的相关信息
	}, cfg.Labels) // 指标的标签名
	prom.MustRegister(vec)
	// 存储到 自定义的结构体 自己在此基础上进行封装
	cv := &promCounterVec{
		counter: vec,
	}

	return cv
}

func (cv *promCounterVec) Inc(labels ...string) {
	cv.counter.WithLabelValues(labels...).Inc()
}

func (cv *promCounterVec) Add(v float64, labels ...string) {
	cv.counter.WithLabelValues(labels...).Add(v)
}
