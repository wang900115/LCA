package bootstrap

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

type promethusManager struct {
	websocketConnections prometheus.Gauge
	requestDuration      *prometheus.HistogramVec
}

type promethusOption struct {
	Namespace string
	Subsystem string
	Enabled   bool
	Buckets   []float64
}

func toFloat64Slice(val interface{}) []float64 {
	switch v := val.(type) {
	case []interface{}:
		result := make([]float64, 0, len(v))
		for _, item := range v {
			switch f := item.(type) {
			case float64:
				result = append(result, f)
			case float32:
				result = append(result, float64(f))
			case int:
				result = append(result, float64(f))
			case int64:
				result = append(result, float64(f))
			case int32:
				result = append(result, float64(f))
			}
		}
		return result
	case []float64:
		return v
	}
	return nil
}

func NewPromethusOption(conf *viper.Viper) promethusOption {
	return promethusOption{
		Namespace: conf.GetString("promethus.name_space"),
		Subsystem: conf.GetString("promethus.subsystem"),
		Enabled:   conf.GetBool("promethus.enabled"),
		Buckets:   toFloat64Slice(conf.Get("promethus.buckets")),
	}
}

func NewPromethus(option promethusOption) *promethusManager {
	if !option.Enabled {
		return nil
	}

	buckets := option.Buckets
	if len(buckets) == 0 {
		buckets = prometheus.DefBuckets
	}

	manager := &promethusManager{
		websocketConnections: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: option.Namespace,
			Subsystem: option.Subsystem,
			Name:      "websocket_connections",
			Help:      "Number of active websocket connections",
		}),

		requestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: option.Namespace,
			Subsystem: option.Subsystem,
			Name:      "http_request_duration_seconds",
			Help:      "Histogram of HTTP request duration",
			Buckets:   buckets,
		}, []string{"path", "method", "status"}),
	}
	prometheus.MustRegister(manager.websocketConnections, manager.requestDuration)
	return manager
}

func (m *promethusManager) MetricsHandler() http.Handler {
	return promhttp.Handler()
}
