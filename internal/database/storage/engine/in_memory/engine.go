package inmemory

import (
	"sync"

	"github.com/crunchydeer30/key-value-database/internal/database/storage/engine"
	"go.uber.org/zap"
)

type InMemoryEngine struct {
	logger *zap.Logger
	mtx    sync.RWMutex
	store  map[string]string
}

func NewInMemoryEngine(logger *zap.Logger) (engine.Engine, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &InMemoryEngine{
		store:  make(map[string]string),
		logger: logger,
		mtx:    sync.RWMutex{},
	}, nil
}

func (e *InMemoryEngine) Get(key string) (string, error) {
	e.logger.Debug("in-memory engine received get command", zap.String("key", key))

	e.mtx.RLock()
	val, ok := e.store[key]
	defer e.mtx.RUnlock()

	e.logger.Debug("Got value", zap.String("value", val))

	if !ok {
		return "", engine.ErrKeyNotFound
	}

	return val, nil
}

func (e *InMemoryEngine) Set(key, value string) error {
	e.logger.Debug(
		"in-memory engine received set command",
		zap.String("key", key),
		zap.String("value", value),
	)

	e.mtx.Lock()
	defer e.mtx.Unlock()

	e.store[key] = value

	return nil
}

func (e *InMemoryEngine) Del(key string) error {
	e.logger.Debug("in-memory engine received delete command", zap.String("key", key))

	e.mtx.Lock()
	defer e.mtx.Unlock()

	delete(e.store, key)

	return nil
}
