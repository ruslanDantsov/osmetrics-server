package repository

import (
	"encoding/json"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum"
	"go.uber.org/zap"
	"os"
	"time"
)

type PersistentStorage struct {
	base          *MemStorage
	filePath      string
	logger        zap.Logger
	ticker        *time.Ticker
	quit          chan struct{}
	storeInterval time.Duration
	isRestore     bool
}

func NewPersistentStorage(base *MemStorage, filePath string, storeInterval time.Duration, logger zap.Logger, isRestore bool) *PersistentStorage {
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

func (ps *PersistentStorage) SaveMetric(m *model.Metrics) (*model.Metrics, error) {
	result, err := ps.base.SaveMetric(m)
	if err != nil {
		return nil, err
	}

	if ps.storeInterval == 0 {
		ps.saveToFile()
	}

	return result, nil
}

func (ps *PersistentStorage) GetMetric(metricID enum.MetricID) (*model.Metrics, bool) {
	return ps.base.GetMetric(metricID)
}

func (ps *PersistentStorage) GetKnownMetrics() []string {
	return ps.base.GetKnownMetrics()
}

func (ps *PersistentStorage) saveToFile() {
	ps.base.mu.RLock()
	defer ps.base.mu.RUnlock()

	file, err := os.Create(ps.filePath)
	if err != nil {
		ps.logger.Error("Failed to create save file", zap.Error(err))
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(ps.base.Storage); err != nil {
		ps.logger.Error("Failed to encode metrics", zap.Error(err))
	}
	ps.logger.Info("Metrics data have been stored")
}

func (ps *PersistentStorage) loadFromFile() {
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

	ps.base.mu.Lock()
	defer ps.base.mu.Unlock()
	ps.base.Storage = data
	ps.logger.Info("Metrics restored from file", zap.Int("count", len(data)))
}
