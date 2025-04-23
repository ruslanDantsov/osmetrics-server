package enum

import (
	"fmt"
	"strings"
)

type MetricId string

const (
	Alloc         MetricId = "Alloc"
	BuckHashSys   MetricId = "BuckHashSys"
	Frees         MetricId = "Frees"
	GCCPUFraction MetricId = "GCCPUFraction"
	GCSys         MetricId = "GCSys"
	HeapAlloc     MetricId = "HeapAlloc"
	HeapIdle      MetricId = "HeapIdle"
	HeapInuse     MetricId = "HeapInuse"
	HeapObjects   MetricId = "HeapObjects"
	HeapReleased  MetricId = "HeapReleased"
	HeapSys       MetricId = "HeapSys"
	LastGC        MetricId = "LastGC"
	Lookups       MetricId = "Lookups"
	MCacheInuse   MetricId = "MCacheInuse"
	MCacheSys     MetricId = "MCacheSys"
	MSpanInuse    MetricId = "MSpanInuse"
	MSpanSys      MetricId = "MSpanSys"
	Mallocs       MetricId = "Mallocs"
	NextGC        MetricId = "NextGC"
	NumForcedGC   MetricId = "NumForcedGC"
	NumGC         MetricId = "NumGC"
	OtherSys      MetricId = "OtherSys"
	PauseTotalNs  MetricId = "PauseTotalNs"
	StackInuse    MetricId = "StackInuse"
	StackSys      MetricId = "StackSys"
	Sys           MetricId = "Sys"
	TotalAlloc    MetricId = "TotalAlloc"
	PollCount     MetricId = "PollCount"
	RandomValue   MetricId = "RandomValue"
)

var validMetricIds = map[MetricId]struct{}{
	Alloc:         {},
	BuckHashSys:   {},
	Frees:         {},
	GCCPUFraction: {},
	GCSys:         {},
	HeapAlloc:     {},
	HeapIdle:      {},
	HeapInuse:     {},
	HeapObjects:   {},
	HeapReleased:  {},
	HeapSys:       {},
	LastGC:        {},
	Lookups:       {},
	MCacheInuse:   {},
	MCacheSys:     {},
	MSpanInuse:    {},
	MSpanSys:      {},
	Mallocs:       {},
	NextGC:        {},
	NumForcedGC:   {},
	NumGC:         {},
	OtherSys:      {},
	PauseTotalNs:  {},
	StackInuse:    {},
	StackSys:      {},
	Sys:           {},
	TotalAlloc:    {},
	PollCount:     {},
	RandomValue:   {},
}

func (m MetricId) IsValid() bool {
	_, ok := validMetricIds[m]
	return ok
}

func (m *MetricId) UnmarshalText(data []byte) error {
	id, err := ParseMetricId(string(data))
	if err != nil {
		return err
	}
	*m = id
	return nil
}

func (m MetricId) MarshalText() ([]byte, error) {
	return []byte(m), nil
}

func (m MetricId) String() string {
	return string(m)
}

func ParseMetricId(s string) (MetricId, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return "", fmt.Errorf("metric ID cannot be empty")
	}

	id := MetricId(trimmed)
	if !id.IsValid() {
		return "", fmt.Errorf("invalid MetricID: %q", trimmed)
	}

	return id, nil
}
