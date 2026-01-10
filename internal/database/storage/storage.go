package storage

import (
	"github.com/crunchydeer30/key-value-database/internal/database/storage/engine"
	"go.uber.org/zap"
)

type Storage struct {
	engine engine.Engine
	logger *zap.Logger
}

func NewStorage(e engine.Engine, logger *zap.Logger) (*Storage, error) {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Storage{
		engine: e,
		logger: logger,
	}, nil
}

func (s *Storage) Get(key string) (string, error) {
	s.logger.Debug("storage received get query", zap.String("key", key))

	return s.engine.Get(key)
}

func (s *Storage) Set(key, value string) error {
	s.logger.Debug("storage received set query", zap.String("key", key), zap.String("value", value))
	return s.engine.Set(key, value)
}

func (s *Storage) Del(key string) error {
	s.logger.Debug("storage received delete query", zap.String("key", key))
	return s.engine.Del(key)
}
