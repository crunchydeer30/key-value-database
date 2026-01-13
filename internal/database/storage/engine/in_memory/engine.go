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
	e.mtx.RLock()
	val, ok := e.store[key]
	defer e.mtx.RUnlock()

	if !ok {
		return "", engine.ErrKeyNotFound
	}

	return val, nil
}

func (e *InMemoryEngine) Set(key, value string) error {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	e.store[key] = value

	return nil
}

func (e *InMemoryEngine) Del(key string) error {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	delete(e.store, key)

	return nil
}
