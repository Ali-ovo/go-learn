package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type (
	// GaugeVecOpts is an alias of VectorOpts.
	GaugeVecOpts VectorOpts

	// GaugeVec represents a gauge vector.
	GaugeVec interface {
		// Set sets v to labels.
		Set(v float64, labels ...string)
		// Inc increments labels.
		Inc(labels ...string)
		// Add adds v to labels.
		Add(v float64, labels ...string)
	}

	promGaugeVec struct {
		// 用于记录可增减的度量数据，即Gauge类型指标。
		// Gauge指标通常用于表示当前状态的度量数据，例如CPU使用率、内存使用量、并发连接数等。
		// GaugeVec可以包含多个维度，每个维度对应一个标签，方便对度量数据进行细粒度的区分和查询。
		gauge *prom.GaugeVec
	}
)

// NewGaugeVec returns a GaugeVec.
func NewGaugeVec(cfg *GaugeVecOpts) GaugeVec {
	if cfg == nil {
		return nil
	}

	vec := prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
		}, cfg.Labels)
	prom.MustRegister(vec)
	gv := &promGaugeVec{
		gauge: vec,
	}

	return gv
}

func (gv *promGaugeVec) Inc(labels ...string) {
	gv.gauge.WithLabelValues(labels...).Inc()
}

func (gv *promGaugeVec) Add(v float64, labels ...string) {
	gv.gauge.WithLabelValues(labels...).Add(v)
}

func (gv *promGaugeVec) Set(v float64, labels ...string) {
	gv.gauge.WithLabelValues(labels...).Set(v)
}
