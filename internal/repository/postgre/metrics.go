package postgre

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/ruslanDantsov/osmetrics-server/internal/constants"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum"
	"github.com/ruslanDantsov/osmetrics-server/internal/repository/postgre/sqlqueries"
	"go.uber.org/zap"
)

type PostgreStorage struct {
	conn *pgxpool.Pool
	Log  zap.Logger
}

func NewPostgreStorage(log zap.Logger, connectionString string) (*PostgreStorage, error) {
	if err := applyMigrations(connectionString); err != nil {
		return nil, err
	}

	conn, err := pgxpool.New(context.Background(), connectionString)

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	return &PostgreStorage{
		conn: conn,
		Log:  log,
	}, nil

}

func (s *PostgreStorage) GetMetric(ctx context.Context, metricID enum.MetricID) (*model.Metrics, bool) {
	var (
		existingID    string
		existingType  string
		existingDelta sql.NullInt64
		existingValue sql.NullFloat64
	)

	err := s.conn.QueryRow(
		ctx,
		sqlqueries.SelectMetricById,
		metricID).
		Scan(&existingID, &existingType, &existingDelta, &existingValue)

	if err != nil {
		return nil, false
	}

	enumExistingId, _ := enum.ParseMetricID(existingID)

	metric := &model.Metrics{
		ID:    enumExistingId,
		MType: existingType,
	}

	if existingDelta.Valid {
		val := existingDelta.Int64
		metric.Delta = &val
	}
	if existingValue.Valid {
		val := existingValue.Float64
		metric.Value = &val
	}

	return metric, true
}

func (s *PostgreStorage) GetKnownMetrics(ctx context.Context) []string {
	var metricNames []string

	rows, err := s.conn.Query(ctx, sqlqueries.SelectAllMetricIDs)

	if err != nil {
		return metricNames
	}
	defer rows.Close()

	for rows.Next() {
		var existingID string
		if err := rows.Scan(&existingID); err != nil {
			return metricNames
		}

		metricNames = append(metricNames, existingID)
	}

	return metricNames

}

func (s *PostgreStorage) SaveMetric(ctx context.Context, metric *model.Metrics) (*model.Metrics, error) {
	switch metric.MType {
	case constants.CounterMetricType:
		return s.saveCounterMetric(ctx, metric)
	case constants.GaugeMetricType:
		return s.saveGaugeMetric(ctx, metric)
	default:
		return nil, fmt.Errorf("unsupported metric type: %s", metric.MType)
	}

}

func (s *PostgreStorage) saveCounterMetric(ctx context.Context, metric *model.Metrics) (*model.Metrics, error) {
	if metric.Delta == nil {
		*metric.Delta = 0
	}

	var (
		existingID    string
		existingType  string
		existingDelta sql.NullInt64
	)

	err := s.conn.QueryRow(
		ctx,
		sqlqueries.SelectMetricById,
		metric.ID).
		Scan(&existingID, &existingType, &existingDelta)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		existingDelta.Int64 = 0
	case err != nil:
		return nil, err
	}

	*metric.Delta += existingDelta.Int64

	_, err = s.conn.Exec(ctx,
		sqlqueries.InsertOrUpdateCounterMetric,
		metric.ID,
		metric.MType,
		metric.Delta)

	if err != nil {
		return nil, err
	}

	return metric, nil
}

func (s *PostgreStorage) saveGaugeMetric(ctx context.Context, metric *model.Metrics) (*model.Metrics, error) {
	if metric.Value == nil {
		*metric.Value = float64(0)
	}

	_, err := s.conn.Exec(ctx,
		sqlqueries.InsertOrUpdateGaugeMetric,
		metric.ID,
		metric.MType,
		*metric.Value)

	if err != nil {
		return nil, err
	}

	return metric, nil
}

func (s *PostgreStorage) HealthCheck(ctx context.Context) error {

	if err := s.conn.Ping(ctx); err != nil {
		s.Log.Warn("DB health check failed", zap.Error(err))
		return err
	}

	s.Log.Info("DB health check passed")
	return nil
}

func (s *PostgreStorage) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func applyMigrations(connectionString string) error {
	sqlDB, err := sql.Open("pgx", connectionString)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(sqlDB, "internal/db/migrations")
}
