package model

import (
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum/metric"
	"strconv"
)

type GaugeMetricModel struct {
	MetricType metric.MetricType `json:"metric_type"`
	Value      float64           `json:"value"`
}

func NewGaugeMetricModelWithRawTypes(metricTypeRaw string, valueRaw string) (*GaugeMetricModel, error) {

	metricType, err := metric.ParseMetricType(metricTypeRaw)
	if err != nil {
		return nil, err
	}

	floatValue, err := strconv.ParseFloat(valueRaw, 64)

	if err != nil {
		return nil, err
	}

	return &GaugeMetricModel{
		MetricType: metricType,
		Value:      floatValue,
	}, nil
}
