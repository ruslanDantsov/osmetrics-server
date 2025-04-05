package metric

import (
	"fmt"
	"strings"
)

type MetricType string

const (
	Counter MetricType = "Counter"
	Gauge   MetricType = "Gauge"
)

type Metric struct {
	Name string
	Type MetricType
}

var (
	Alloc = Metric{
		Name: "Alloc",
		Type: Gauge,
	}
	BuckHashSys = Metric{
		Name: "BuckHashSys",
		Type: Gauge,
	}

	Frees = Metric{
		Name: "Frees",
		Type: Gauge,
	}
	GCCPUFraction = Metric{
		Name: "GCCPUFraction",
		Type: Gauge,
	}
	GCSys = Metric{
		Name: "GCSys",
		Type: Gauge,
	}
	HeapAlloc = Metric{
		Name: "HeapAlloc",
		Type: Gauge,
	}
	HeapIdle = Metric{
		Name: "HeapIdle",
		Type: Gauge,
	}
	HeapInuse = Metric{
		Name: "HeapInuse",
		Type: Gauge,
	}
	HeapObjects = Metric{
		Name: "HeapObjects",
		Type: Gauge,
	}
	HeapReleased = Metric{
		Name: "HeapReleased",
		Type: Gauge,
	}
	HeapSys = Metric{
		Name: "HeapSys",
		Type: Gauge,
	}
	LastGC = Metric{
		Name: "LastGC",
		Type: Gauge,
	}
	Lookups = Metric{
		Name: "Lookups",
		Type: Gauge,
	}
	MCacheInuse = Metric{
		Name: "MCacheInuse",
		Type: Gauge,
	}
	MCacheSys = Metric{
		Name: "MCacheSys",
		Type: Gauge,
	}
	MSpanInuse = Metric{
		Name: "MSpanInuse",
		Type: Gauge,
	}
	MSpanSys = Metric{
		Name: "MSpanSys",
		Type: Gauge,
	}
	Mallocs = Metric{
		Name: "Mallocs",
		Type: Gauge,
	}
	NextGC = Metric{
		Name: "NextGC",
		Type: Gauge,
	}
	NumForcedGC = Metric{
		Name: "NumForcedGC",
		Type: Gauge,
	}
	NumGC = Metric{
		Name: "NumGC",
		Type: Gauge,
	}
	OtherSys = Metric{
		Name: "OtherSys",
		Type: Gauge,
	}
	PauseTotalNs = Metric{
		Name: "PauseTotalNs",
		Type: Gauge,
	}
	StackInuse = Metric{
		Name: "StackInuse",
		Type: Gauge,
	}
	StackSys = Metric{
		Name: "StackSys",
		Type: Gauge,
	}
	Sys = Metric{
		Name: "Sys",
		Type: Gauge,
	}
	TotalAlloc = Metric{
		Name: "TotalAlloc",
		Type: Gauge,
	}
	PollCount = Metric{
		Name: "PollCount",
		Type: Counter,
	}
	RandomValue = Metric{
		Name: "RandomValue",
		Type: Gauge,
	}
)

// var metricNames = []Metric{
// Alloc,
// BuckHashSys,
// Frees,
// GCCPUFraction,
// }
//
// var metricNameMap = func() map[string]Metric {
// m := make(map[string]Metric)
// for _, mt := range metricNames {
// m[string(mt)] = mt
// }
// return m
// }()
func (m Metric) String() string {
	return m.Name
}

//
//func ListMetricNames() []Metric {
//return metricNames
//}

func ParseMetricName(s string) (Metric, error) {
	//TODO: uncomment in  iter2
	//if mt, exists := metricNameMap[s]; exists {
	//	return mt, nil
	//}
	//return "", fmt.Errorf("invalid Metric: %s", s)
	if strings.TrimSpace(s) == "" {
		return Metric{}, fmt.Errorf("invalid Metric: %s", s)
	}

	return Metric{
		Name: s,
		Type: Gauge,
	}, nil
}
