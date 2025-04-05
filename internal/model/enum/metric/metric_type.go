package metric

import (
	"fmt"
)

type MetricType string

const (
	Alloc         MetricType = "Alloc"
	BuckHashSys   MetricType = "BuckHashSys"
	Frees         MetricType = "Frees"
	GCCPUFraction MetricType = "GCCPUFraction"
)

var metricTypes = []MetricType{
	Alloc,
	BuckHashSys,
	Frees,
	GCCPUFraction,
}

var metricTypeMap = func() map[string]MetricType {
	m := make(map[string]MetricType)
	for _, mt := range metricTypes {
		m[string(mt)] = mt
	}
	return m
}()

func (m MetricType) String() string {
	return string(m)
}

func ListMetricTypes() []MetricType {
	return metricTypes
}

func ParseMetricType(s string) (MetricType, error) {
	if mt, exists := metricTypeMap[s]; exists {
		return mt, nil
	}
	return "", fmt.Errorf("invalid MetricType: %s", s)
}
