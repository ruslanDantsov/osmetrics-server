package constants

import "time"

const (
	URLParamMetricType  = "type"
	URLParamMetricName  = "name"
	URLParamMetricValue = "value"
	GaugeMetricType     = "gauge"
	CounterMetricType   = "counter"
	DBPingTimeout       = 1 * time.Second
)
