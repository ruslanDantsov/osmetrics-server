package constants

import "time"

const (
	URLParamMetricType        = "type"
	URLParamMetricName        = "name"
	URLParamMetricValue       = "value"
	GaugeMetricType           = "gauge"
	CounterMetricType         = "counter"
	DBPingTimeout             = 10 * time.Second
	DBRetryConnectMaxAttempts = 10
	DBRetryDelayTime          = 3 * time.Second
)
