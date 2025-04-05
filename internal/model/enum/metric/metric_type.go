package metric

import (
	"fmt"
	"strings"
)

type MetricName string

const (
	Alloc         MetricName = "Alloc"
	BuckHashSys   MetricName = "BuckHashSys"
	Frees         MetricName = "Frees"
	GCCPUFraction MetricName = "GCCPUFraction"
)

var metricNames = []MetricName{
	Alloc,
	BuckHashSys,
	Frees,
	GCCPUFraction,
}

var metricNameMap = func() map[string]MetricName {
	m := make(map[string]MetricName)
	for _, mt := range metricNames {
		m[string(mt)] = mt
	}
	return m
}()

func (m MetricName) String() string {
	return string(m)
}

func ListMetricNames() []MetricName {
	return metricNames
}

func ParseMetricName(s string) (MetricName, error) {
	//TODO: uncomment in  iter2
	//if mt, exists := metricNameMap[s]; exists {
	//	return mt, nil
	//}
	//return "", fmt.Errorf("invalid MetricName: %s", s)
	if strings.TrimSpace(s) == "" {
		return "", fmt.Errorf("invalid MetricName: %s", s)
	}
	return MetricName(s), nil
}
