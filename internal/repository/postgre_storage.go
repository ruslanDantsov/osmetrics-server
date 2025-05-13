package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ruslanDantsov/osmetrics-server/internal/constants"
	"go.uber.org/zap"
	"time"
)

type PostgreStorage struct {
	conn *pgxpool.Pool
	Log  zap.Logger
}

func NewPostgreStorage(log zap.Logger, connectionString string) (*PostgreStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBRetryConnectMaxAttempts*constants.DBRetryDelayTime)
	defer cancel()

	conn, err := connectWithRetries(ctx, connectionString, log)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	return &PostgreStorage{
		conn: conn,
		Log:  log,
	}, nil

}

func connectWithRetries(ctx context.Context, connectionString string, log zap.Logger) (*pgxpool.Pool, error) {
	var conn *pgxpool.Pool
	var err error

	for attempt := 1; attempt <= constants.DBRetryConnectMaxAttempts; attempt++ {
		conn, err = pgxpool.New(context.Background(), connectionString)
		if err == nil {
			if err = conn.Ping(ctx); err == nil {
				return conn, nil
			}
			conn.Close()
		}

		if attempt < constants.DBRetryConnectMaxAttempts {
			log.Info("Database connection failed, retrying...",
				zap.Int("attempt", attempt),
				zap.Int("max_attempts", constants.DBRetryConnectMaxAttempts),
				zap.Error(err),
			)
			time.Sleep(constants.DBRetryDelayTime)
		}
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", constants.DBRetryConnectMaxAttempts, err)
}

func (s *PostgreStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBPingTimeout)
	defer cancel()

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
