package sqlqueries

const (
	SelectAllMetricIDs = `
	SELECT id FROM metrics;
`

	InsertOrUpdateGaugeMetric = `
		INSERT INTO metrics (id, type, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET value = $3;
	`

	SelectMetricById = `
	SELECT id, type, delta, value FROM metrics
	WHERE id = $1;
`

	InsertOrUpdateCounterMetric = `
		INSERT INTO metrics (id, type, delta)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET delta = $3;
	`
)
