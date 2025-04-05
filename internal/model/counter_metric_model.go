package model

import (
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum/metric"
	"strconv"
)

type CounterMetricModel struct {
	MetricType metric.MetricType
	Value      int64
}

func NewCounterMetricModelWithRawTypes(metricTypeRaw string, valueRaw string) (*CounterMetricModel, error) {
	metricType, err := metric.ParseMetricType(metricTypeRaw)
	if err != nil {
		return nil, err
	}

	intValue, err := strconv.ParseInt(valueRaw, 10, 64)

	if err != nil {
		return nil, err
	}

	return &CounterMetricModel{
		MetricType: metricType,
		Value:      intValue,
	}, nil
}
