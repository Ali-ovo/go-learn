package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type (
	// A HistogramVecOpts is a histogram vector options.
	HistogramVecOpts struct {
		Namespace string
		Subsystem string
		Name      string
		Help      string
		Labels    []string
		Buckets   []float64
	}

	// A HistogramVec interface represents a histogram vector.
	HistogramVec interface {
		// Observe adds observation v to labels.
		Observe(v int64, labels ...string)
	}

	promHistogramVec struct {
		// 用于记录分布式度量数据，即Histogram类型指标。
		// Histogram指标通常用于表示度量数据的分布情况，例如请求响应时间、API调用次数、文件大小等。
		// HistogramVec同样可以包含多个维度，每个维度对应一个标签，方便对分布式数据进行细粒度的区分和查询。
		histogram *prom.HistogramVec
	}
)

// NewHistogramVec returns a HistogramVec.
func NewHistogramVec(cfg *HistogramVecOpts) HistogramVec {
	if cfg == nil {
		return nil
	}

	vec := prom.NewHistogramVec(prom.HistogramOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      cfg.Name,
		Help:      cfg.Help,
		Buckets:   cfg.Buckets,
	}, cfg.Labels)
	prom.MustRegister(vec)
	hv := &promHistogramVec{
		histogram: vec,
	}

	return hv
}

func (hv *promHistogramVec) Observe(v int64, labels ...string) {
	hv.histogram.WithLabelValues(labels...).Observe(float64(v))
}
