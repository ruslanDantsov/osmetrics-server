// Package postgre provides persistent storage implementation that saves metrics data to DB
package postgre

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	if err := applyMigrations(connectionString, log); err != nil {
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
		sqlqueries.SelectMetricByID,
		metricID).
		Scan(&existingID, &existingType, &existingDelta, &existingValue)

	if err != nil {
		return nil, false
	}

	enumExistingID, _ := enum.ParseMetricID(existingID)

	metric := &model.Metrics{
		ID:    enumExistingID,
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

func (s *PostgreStorage) SaveAllMetrics(ctx context.Context, metricList model.MetricsList) (model.MetricsList, error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			s.Log.Error(err.Error())
		}
	}()

	for _, metric := range metricList {
		switch metric.MType {
		case constants.CounterMetricType:
			err = s.saveCounterMetricTx(ctx, tx, &metric)
		case constants.GaugeMetricType:
			err = s.saveGaugeMetricTx(ctx, tx, &metric)
		default:
			err = fmt.Errorf("unsupported metric type: %s", metric.MType)
		}

		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return metricList, nil
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
		zero := int64(0)
		metric.Delta = &zero
	}

	var existingDelta sql.NullInt64

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			s.Log.Error(err.Error())
		}
	}()

	err = tx.QueryRow(
		ctx,
		sqlqueries.SelectMetricByID,
		metric.ID).
		Scan(new(string), new(string), &existingDelta, new(sql.NullFloat64))

	switch {
	case errors.Is(err, sql.ErrNoRows):
		existingDelta.Int64 = 0
	case err != nil:
		return nil, err
	}

	*metric.Delta += existingDelta.Int64

	_, err = tx.Exec(ctx,
		sqlqueries.InsertOrUpdateCounterMetric,
		metric.ID,
		metric.MType,
		metric.Delta)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return nil, fmt.Errorf("unique constraint violation when saving metric: %w", err)
			default:
				return nil, fmt.Errorf("postgresql error when saving metric (code %s): %w", pgErr.Code, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return metric, nil
}

func (s *PostgreStorage) saveCounterMetricTx(ctx context.Context, tx pgx.Tx, metric *model.Metrics) error {
	if metric.Delta == nil {
		zero := int64(0)
		metric.Delta = &zero
	}

	var existingDelta sql.NullInt64
	err := tx.QueryRow(
		ctx,
		sqlqueries.SelectMetricByID,
		metric.ID,
	).Scan(new(string), new(string), &existingDelta, new(sql.NullFloat64))

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("select counter metric failed: %w", err)
	}

	if existingDelta.Valid {
		*metric.Delta += existingDelta.Int64
	}

	_, err = tx.Exec(ctx,
		sqlqueries.InsertOrUpdateCounterMetric,
		metric.ID,
		metric.MType,
		metric.Delta)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return fmt.Errorf("unique constraint violation when saving metric: %w", err)
			default:
				return fmt.Errorf("postgresql error when saving metric (code %s): %w", pgErr.Code, err)
			}
		}
	}

	return nil
}

func (s *PostgreStorage) saveGaugeMetric(ctx context.Context, metric *model.Metrics) (*model.Metrics, error) {
	if metric.Value == nil {
		zero := float64(0)
		metric.Value = &zero
	}

	_, err := s.conn.Exec(ctx,
		sqlqueries.InsertOrUpdateGaugeMetric,
		metric.ID,
		metric.MType,
		*metric.Value)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return nil, fmt.Errorf("unique constraint violation when saving metric: %w", err)
			default:
				return nil, fmt.Errorf("postgresql error when saving metric (code %s): %w", pgErr.Code, err)
			}
		}
	}

	return metric, nil
}

func (s *PostgreStorage) saveGaugeMetricTx(ctx context.Context, tx pgx.Tx, metric *model.Metrics) error {
	if metric.Value == nil {
		zero := float64(0)
		metric.Value = &zero
	}

	_, err := tx.Exec(ctx,
		sqlqueries.InsertOrUpdateGaugeMetric,
		metric.ID,
		metric.MType,
		*metric.Value)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return fmt.Errorf("unique constraint violation when saving metric: %w", err)
			default:
				return fmt.Errorf("postgresql error when saving metric (code %s): %w", pgErr.Code, err)
			}
		}
	}

	return nil
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

func applyMigrations(connectionString string, log zap.Logger) error {
	sqlDB, err := sql.Open("pgx", connectionString)
	if err != nil {
		return err
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Error(err.Error())
		}
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(sqlDB, "internal/db/migrations")
}
