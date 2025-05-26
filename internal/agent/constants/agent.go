package constants

import "time"

const (
	MaxDelayForWaitingServer      = 10 * time.Second
	IncreaseDelayForWaitingServer = 2 * time.Second
	ServerHealthCheckURL          = "http://%v/health"
	HashHeaderName                = "HashSHA256"
	MetricChannelSize             = 100
	UpdateMetricURL               = "http://%v/update"
	UpdateMetricsURL              = "http://%v/updates"
	GaugeMetricType               = "gauge"
	CounterMetricType             = "counter"
)
