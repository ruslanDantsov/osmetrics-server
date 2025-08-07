package constants

const (
	// URLParamMetricType — имя параметра URL, указывающее тип метрики (например, "gauge" или "counter").
	URLParamMetricType = "type"

	// URLParamMetricName — имя параметра URL, указывающее название метрики.
	URLParamMetricName = "name"

	// URLParamMetricValue — имя параметра URL, указывающее значение метрики.
	URLParamMetricValue = "value"

	// GaugeMetricType — тип метрики "gauge"
	GaugeMetricType = "gauge"

	// CounterMetricType — тип метрики "counter"
	CounterMetricType = "counter"

	// HashHeaderName — имя HTTP-заголовка, содержащего хеш-сумму (SHA256) тела запроса для проверки целостности данных.
	HashHeaderName = "HashSHA256"
)
