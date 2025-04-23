package model

import (
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum"
	"strconv"
)

//go:generate easyjson -all Metrics.go
type Metrics struct {
	ID    enum.MetricID `json:"id"`
	MType string        `json:"type"`
	Delta *int64        `json:"delta,omitempty"`
	Value *float64      `json:"value,omitempty"`
}

func NewMetricWithRawValues(metricType string, metricIdRaw string, valueRaw string) (*Metrics, error) {
	metricId, err := enum.ParseMetricID(metricIdRaw)
	if err != nil {
		return nil, err
	}

	switch metricType {
	case "Gauge":
		floatValue, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			return nil, err
		}
		return &Metrics{
			ID:    metricId,
			MType: metricType,
			Value: &floatValue,
		}, nil

	case "Counter":
		intValue, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil {
			return nil, err
		}
		return &Metrics{
			ID:    metricId,
			MType: metricType,
			Delta: &intValue,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported metric type: %s", metricType)
	}
}
