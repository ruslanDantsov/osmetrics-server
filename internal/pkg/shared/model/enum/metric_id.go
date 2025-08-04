package enum

import (
	"fmt"
	"strings"
)

// MetricID представляет собой строковый идентификатор метрики.
type MetricID string

// Список возможных значений идентификаторов метрик.
const (
	Alloc           MetricID = "Alloc"
	BuckHashSys     MetricID = "BuckHashSys"
	Frees           MetricID = "Frees"
	GCCPUFraction   MetricID = "GCCPUFraction"
	GCSys           MetricID = "GCSys"
	HeapAlloc       MetricID = "HeapAlloc"
	HeapIdle        MetricID = "HeapIdle"
	HeapInuse       MetricID = "HeapInuse"
	HeapObjects     MetricID = "HeapObjects"
	HeapReleased    MetricID = "HeapReleased"
	HeapSys         MetricID = "HeapSys"
	LastGC          MetricID = "LastGC"
	Lookups         MetricID = "Lookups"
	MCacheInuse     MetricID = "MCacheInuse"
	MCacheSys       MetricID = "MCacheSys"
	MSpanInuse      MetricID = "MSpanInuse"
	MSpanSys        MetricID = "MSpanSys"
	Mallocs         MetricID = "Mallocs"
	NextGC          MetricID = "NextGC"
	NumForcedGC     MetricID = "NumForcedGC"
	NumGC           MetricID = "NumGC"
	OtherSys        MetricID = "OtherSys"
	PauseTotalNs    MetricID = "PauseTotalNs"
	StackInuse      MetricID = "StackInuse"
	StackSys        MetricID = "StackSys"
	Sys             MetricID = "Sys"
	TotalAlloc      MetricID = "TotalAlloc"
	PollCount       MetricID = "PollCount"
	RandomValue     MetricID = "RandomValue"
	TotalMemory     MetricID = "TotalMemory"
	FreeMemory      MetricID = "FreeMemory"
	CPUutilization1 MetricID = "CPUutilization1"
)

var validMetricIDs = map[MetricID]struct{}{
	Alloc:           {},
	BuckHashSys:     {},
	Frees:           {},
	GCCPUFraction:   {},
	GCSys:           {},
	HeapAlloc:       {},
	HeapIdle:        {},
	HeapInuse:       {},
	HeapObjects:     {},
	HeapReleased:    {},
	HeapSys:         {},
	LastGC:          {},
	Lookups:         {},
	MCacheInuse:     {},
	MCacheSys:       {},
	MSpanInuse:      {},
	MSpanSys:        {},
	Mallocs:         {},
	NextGC:          {},
	NumForcedGC:     {},
	NumGC:           {},
	OtherSys:        {},
	PauseTotalNs:    {},
	StackInuse:      {},
	StackSys:        {},
	Sys:             {},
	TotalAlloc:      {},
	PollCount:       {},
	RandomValue:     {},
	TotalMemory:     {},
	FreeMemory:      {},
	CPUutilization1: {},
}

// IsValid проверяет, является ли идентификатор метрики допустимым.
func (m *MetricID) IsValid() bool {
	_, ok := validMetricIDs[*m]
	return ok
}

// UnmarshalText преобразует текст в MetricID и проверяет его валидность.
func (m *MetricID) UnmarshalText(data []byte) error {
	id, err := ParseMetricID(string(data))
	if err != nil {
		return err
	}
	*m = id
	return nil
}

// MarshalText преобразует MetricID в срез байт.
func (m *MetricID) MarshalText() ([]byte, error) {
	return []byte(*m), nil
}

// String возвращает строковое представление идентификатора метрики.
func (m *MetricID) String() string {
	return string(*m)
}

// ParseMetricID возвращает MetricID, созданный из переданной строки,
// или ошибку, если строка пустая.
func ParseMetricID(s string) (MetricID, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return "", fmt.Errorf("metric ID cannot be empty")
	}

	id := MetricID(trimmed)
	return id, nil
}
