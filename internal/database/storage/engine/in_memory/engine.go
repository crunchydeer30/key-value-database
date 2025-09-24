package inmemory

import (
	"kv-db/internal/database/storage/engine"

	"go.uber.org/zap"
)

type InMemoryEngine struct {
	logger *zap.Logger
	store  map[string]string
}

func NewInMemoryEngine(logger *zap.Logger) (engine.Engine, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &InMemoryEngine{
		store:  make(map[string]string),
		logger: logger,
	}, nil
}

func (e *InMemoryEngine) Get(key string) (string, error) {
	e.logger.Debug("in-memory engine received get command", zap.String("key", key))

	val, ok := e.store[key]

	e.logger.Debug("Got value", zap.String("value", val))

	if !ok {
		return "", engine.ErrKeyNotFound
	}

	return val, nil
}

func (e *InMemoryEngine) Set(key, value string) error {
	e.logger.Debug("in-memory engine received set command", zap.String("key", key), zap.String("value", value))

	e.store[key] = value

	return nil
}

func (e *InMemoryEngine) Del(key string) error {
	e.logger.Debug("in-memory engine received delete command", zap.String("key", key))

	delete(e.store, key)

	return nil
}
