// Package constants provides constants for the osmetrics server.
package constants

import "time"

const (
	URLParamMetricType  = "type"
	URLParamMetricName  = "name"
	URLParamMetricValue = "value"
	GaugeMetricType     = "gauge"
	CounterMetricType   = "counter"
	DBPingTimeout       = 10 * time.Second
	DBQueryTimeout      = 10 * time.Second
)
