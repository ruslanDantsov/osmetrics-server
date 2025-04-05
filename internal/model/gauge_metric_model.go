package model

import (
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum/metric"
	"strconv"
)

type GaugeMetricModel struct {
	Name  metric.MetricName
	Value float64
}

func NewGaugeMetricModelWithRawValues(metricNameRaw string, valueRaw string) (*GaugeMetricModel, error) {

	metricName, err := metric.ParseMetricName(metricNameRaw)
	if err != nil {
		return nil, err
	}

	floatValue, err := strconv.ParseFloat(valueRaw, 64)

	if err != nil {
		return nil, err
	}

	return &GaugeMetricModel{
		Name:  metricName,
		Value: floatValue,
	}, nil
}
