package model

import (
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum/metric"
	"strconv"
)

type CounterMetricModel struct {
	Name  metric.MetricName
	Value int64
}

func NewCounterMetricModelWithRawValues(metricNameRaw string, valueRaw string) (*CounterMetricModel, error) {
	metricName, err := metric.ParseMetricName(metricNameRaw)
	if err != nil {
		return nil, err
	}

	intValue, err := strconv.ParseInt(valueRaw, 10, 64)

	if err != nil {
		return nil, err
	}

	return &CounterMetricModel{
		Name:  metricName,
		Value: intValue,
	}, nil
}
