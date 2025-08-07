package constants

import "time"

const (
	// MaxDelayForWaitingServer определяет время задержки перед повторной попыткой подключения к серверу.
	MaxDelayForWaitingServer = 10 * time.Second

	// IncreaseDelayForWaitingServer определяет прирост задержки между попытками подключения к серверу.
	IncreaseDelayForWaitingServer = 2 * time.Second

	// ServerHealthCheckURL шаблон URL для проверки состояния сервера.
	ServerHealthCheckURL = "http://%v/health"

	// HashHeaderName — имя HTTP-заголовка, в котором передаётся SHA256-хеш содержимого.
	HashHeaderName = "HashSHA256"

	// MetricChannelSize определяет размер буфера канала метрик.
	MetricChannelSize = 100

	// UpdateMetricURL шаблон URL для обновления одной метрики.
	UpdateMetricURL = "http://%v/update"

	// UpdateMetricsURL шаблон URL для пакетного обновления метрик.
	UpdateMetricsURL = "http://%v/updates"

	// GaugeMetricType указывает тип метрики "gauge"
	GaugeMetricType = "gauge"

	// CounterMetricType указывает тип метрики "counter"
	CounterMetricType = "counter"
)
