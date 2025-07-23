package model

import (
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model/enum"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/constants"
	"strconv"
)

//go:generate easyjson -all Metrics.go

// Metrics структура, которая может быть либо счетчиком (Counter), либо измеряемым значением (Gauge).
type Metrics struct {
	ID    enum.MetricID `json:"id"`              // Уникальный идентификатор метрики
	MType string        `json:"type"`            // Тип метрики: "gauge" или "counter"
	Delta *int64        `json:"delta,omitempty"` // Значение для счетчика (Counter); применяется, если тип метрики — "counter"
	Value *float64      `json:"value,omitempty"` // Значение для измеряемой метрики (Gauge); применяется, если тип метрики — "gauge"
}

// MetricsList представляет собой список метрик.
//
//easyjson:json
type MetricsList []Metrics

// NewMetricWithRawValues создает новую метрику из строковых представлений типа, идентификатора и значения.
func NewMetricWithRawValues(metricType string, metricIDRaw string, valueRaw string) (*Metrics, error) {
	metricID, err := enum.ParseMetricID(metricIDRaw)
	if err != nil {
		return nil, err
	}

	switch metricType {
	case constants.GaugeMetricType:
		floatValue, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			return nil, err
		}
		return &Metrics{
			ID:    metricID,
			MType: metricType,
			Value: &floatValue,
		}, nil

	case constants.CounterMetricType:
		intValue, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil {
			return nil, err
		}
		return &Metrics{
			ID:    metricID,
			MType: metricType,
			Delta: &intValue,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported metric type: %s", metricType)
	}
}
