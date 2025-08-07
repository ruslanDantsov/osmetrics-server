// Package model defines the core data structures of Metrics
package model

import (
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum"
	"strconv"
)

//go:generate easyjson -all Metrics.go

var (
	GaugeMetricType   = "gauge"
	CounterMetricType = "counter"
)

type Metrics struct {
	ID    enum.MetricID `json:"id"`
	MType string        `json:"type"`
	Delta *int64        `json:"delta,omitempty"`
	Value *float64      `json:"value,omitempty"`
}

//easyjson:json
type MetricsList []Metrics

func NewMetricWithRawValues(metricType string, metricIDRaw string, valueRaw string) (*Metrics, error) {
	metricID, err := enum.ParseMetricID(metricIDRaw)
	if err != nil {
		return nil, err
	}

	switch metricType {
	case GaugeMetricType:
		floatValue, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			return nil, err
		}
		return &Metrics{
			ID:    metricID,
			MType: metricType,
			Value: &floatValue,
		}, nil

	case CounterMetricType:
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
