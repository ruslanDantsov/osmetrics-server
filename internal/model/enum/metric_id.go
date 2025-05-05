package enum

import (
	"fmt"
	"strings"
)

type MetricID string

const (
	Alloc         MetricID = "Alloc"
	BuckHashSys   MetricID = "BuckHashSys"
	Frees         MetricID = "Frees"
	GCCPUFraction MetricID = "GCCPUFraction"
	GCSys         MetricID = "GCSys"
	HeapAlloc     MetricID = "HeapAlloc"
	HeapIdle      MetricID = "HeapIdle"
	HeapInuse     MetricID = "HeapInuse"
	HeapObjects   MetricID = "HeapObjects"
	HeapReleased  MetricID = "HeapReleased"
	HeapSys       MetricID = "HeapSys"
	LastGC        MetricID = "LastGC"
	Lookups       MetricID = "Lookups"
	MCacheInuse   MetricID = "MCacheInuse"
	MCacheSys     MetricID = "MCacheSys"
	MSpanInuse    MetricID = "MSpanInuse"
	MSpanSys      MetricID = "MSpanSys"
	Mallocs       MetricID = "Mallocs"
	NextGC        MetricID = "NextGC"
	NumForcedGC   MetricID = "NumForcedGC"
	NumGC         MetricID = "NumGC"
	OtherSys      MetricID = "OtherSys"
	PauseTotalNs  MetricID = "PauseTotalNs"
	StackInuse    MetricID = "StackInuse"
	StackSys      MetricID = "StackSys"
	Sys           MetricID = "Sys"
	TotalAlloc    MetricID = "TotalAlloc"
	PollCount     MetricID = "PollCount"
	RandomValue   MetricID = "RandomValue"
)

var validMetricIDs = map[MetricID]struct{}{
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

func (m MetricID) IsValid() bool {
	_, ok := validMetricIDs[m]
	return ok
}

func (m *MetricID) UnmarshalText(data []byte) error {
	id, err := ParseMetricID(string(data))
	if err != nil {
		return err
	}
	*m = id
	return nil
}

func (m MetricID) MarshalText() ([]byte, error) {
	return []byte(m), nil
}

func (m MetricID) String() string {
	return string(m)
}

func ParseMetricID(s string) (MetricID, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return "", fmt.Errorf("metric ID cannot be empty")
	}

	id := MetricID(trimmed)
	return id, nil
}
