package storage

import (
	"errors"

	"github.com/crunchydeer30/key-value-database/internal/database/storage/engine"
	"go.uber.org/zap"
)

type Storage struct {
	engine engine.Engine
	logger *zap.Logger
}

func NewStorage(e engine.Engine, logger *zap.Logger) (*Storage, error) {
	if e == nil {
		return nil, errors.New("engine is nil")
	}

	if logger == nil {
		logger = zap.NewNop()
	}

	return &Storage{
		engine: e,
		logger: logger,
	}, nil
}

func (s *Storage) Get(key string) (string, error) {
	return s.engine.Get(key)
}

func (s *Storage) Set(key, value string) error {
	return s.engine.Set(key, value)
}

func (s *Storage) Del(key string) error {
	return s.engine.Del(key)
}
