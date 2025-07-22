package file

import (
	"context"
	"encoding/json"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model/enum"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/repository/memory"
	"go.uber.org/zap"
	"os"
	"time"
)

type MemoryStorager interface {
	GetKnownMetrics(ctx context.Context) []string
	GetMetric(ctx context.Context, metricID enum.MetricID) (*model.Metrics, bool)
	SaveMetric(ctx context.Context, metric *model.Metrics) (*model.Metrics, error)
	SaveAllMetrics(ctx context.Context, metricList model.MetricsList) (model.MetricsList, error)
	HealthCheck(ctx context.Context) error
	Close()
}

type PersistentStorage struct {
	base          MemoryStorager
	filePath      string
	logger        zap.Logger
	ticker        *time.Ticker
	quit          chan struct{}
	storeInterval time.Duration
	isRestore     bool
}

func NewPersistentStorage(base MemoryStorager, filePath string, storeInterval time.Duration, logger zap.Logger, isRestore bool) *PersistentStorage {
	ps := &PersistentStorage{
		base:          base,
		filePath:      filePath,
		logger:        logger,
		quit:          make(chan struct{}),
		storeInterval: storeInterval,
		isRestore:     isRestore,
	}
	if isRestore {
		ps.loadFromFile()
	}

	if storeInterval > 0 {
		ps.StartAutoSave(storeInterval)
	}
	return ps
}

func (ps *PersistentStorage) HealthCheck(ctx context.Context) error {
	return nil
}

func (ps *PersistentStorage) Close() {
	//For this type of storage we don't need implementation
}

func (ps *PersistentStorage) StartAutoSave(interval time.Duration) {
	ps.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ps.ticker.C:
				ps.saveToFile()
			case <-ps.quit:
				ps.ticker.Stop()
				ps.saveToFile()
				return
			}
		}
	}()
}

func (ps *PersistentStorage) Shutdown() {
	if ps.storeInterval > 0 {
		close(ps.quit)
	}
	ps.saveToFile()
}

func (ps *PersistentStorage) SaveMetric(ctx context.Context, m *model.Metrics) (*model.Metrics, error) {
	result, err := ps.base.SaveMetric(ctx, m)
	if err != nil {
		return nil, err
	}

	if ps.storeInterval == 0 {
		ps.saveToFile()
	}

	return result, nil
}
func (ps *PersistentStorage) SaveAllMetrics(ctx context.Context, metricList model.MetricsList) (model.MetricsList, error) {
	result, err := ps.base.SaveAllMetrics(ctx, metricList)
	if err != nil {
		return nil, err
	}

	if ps.storeInterval == 0 {
		ps.saveToFile()
	}

	return result, nil
}

func (ps *PersistentStorage) GetMetric(ctx context.Context, metricID enum.MetricID) (*model.Metrics, bool) {
	return ps.base.GetMetric(ctx, metricID)
}

func (ps *PersistentStorage) GetKnownMetrics(ctx context.Context) []string {
	return ps.base.GetKnownMetrics(ctx)
}

func (ps *PersistentStorage) saveToFile() {
	memStorage, ok := ps.base.(*memory.MemStorage)
	if !ok {
		ps.logger.Error("saveToFile: base is not MemStorage, skipping file save")
		return
	}

	memStorage.Mu.RLock()
	defer memStorage.Mu.RUnlock()

	file, err := os.Create(ps.filePath)
	if err != nil {
		ps.logger.Error("Failed to create save file", zap.Error(err))
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(memStorage.Storage); err != nil {
		ps.logger.Error("Failed to encode metrics", zap.Error(err))
	}
	ps.logger.Info("Metrics data have been stored")
}

func (ps *PersistentStorage) loadFromFile() {
	memStorage, ok := ps.base.(*memory.MemStorage)
	if !ok {
		ps.logger.Error("loadFromFile: base is not MemStorage, skipping file load")
		return
	}

	file, err := os.Open(ps.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			ps.logger.Info("No restore file found; starting fresh")
			return
		}
		ps.logger.Error("Failed to open restore file", zap.Error(err))
		return
	}
	defer file.Close()

	data := make(map[string]*model.Metrics)
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		ps.logger.Error("Failed to decode metrics file", zap.Error(err))
		return
	}

	memStorage.Mu.Lock()
	defer memStorage.Mu.Unlock()
	memStorage.Storage = data
	ps.logger.Info("Metrics restored from file", zap.Int("count", len(data)))
}
