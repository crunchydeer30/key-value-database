package database

import (
	"errors"
	"fmt"
	"kv-db/internal/database/compute"
	"kv-db/internal/database/storage"
	"kv-db/internal/database/storage/engine"
	inmemory "kv-db/internal/database/storage/engine/in_memory"

	"go.uber.org/zap"
)

type Database struct {
	compute *compute.Compute
	storage *storage.Storage
	logger  *zap.Logger
}

func NewDatabase(logger *zap.Logger) (*Database, error) {
	if logger == nil {
		return nil, errors.New("no logger is provided")
	}
	logger.Debug("initializing database...")

	logger.Debug("initializing compute layer...")

	compute, err := compute.NewCompute(logger)
	if err != nil {
		logger.Error("failed to initialize compute layer", zap.Error(err))
		return nil, err
	}

	logger.Debug("compute layer initialized")

	logger.Debug("initializing storage layer...")

	engine, err := inmemory.NewInMemoryEngine(logger)
	if err != nil {
		logger.Error("failed to initialize storage layer", zap.Error(err))
		return nil, err
	}

	storage, err := storage.NewStorage(engine, logger)

	if err != nil {
		logger.Error("failed to initialize storage layer", zap.Error(err))
		return nil, err
	}
	logger.Debug("storage layer initialized")

	logger.Debug("database initialized")
	return &Database{
		compute: compute,
		storage: storage,
		logger:  logger,
	}, nil
}

func (d *Database) HandleQuery(queryStr string) string {
	d.logger.Debug("database received query", zap.String("query", queryStr))

	query, err := d.compute.Parser.Parse(queryStr)

	if err != nil {
		d.logger.Debug("invalid query", zap.String("query", queryStr), zap.Error(err))
		return fmt.Sprintf("invalid query: %s", err.Error())
	}

	switch query.Command {
	case compute.CommandName("GET"):
		return d.handleGetQuery(query)
	case compute.CommandName("SET"):
		return d.handleSetQuery(query)
	case compute.CommandName("DEL"):
		return d.handleDelQuery(query)
	}

	return "internal error"
}

func (d *Database) handleGetQuery(query *compute.Query) string {
	val, err := d.storage.Get(query.Args[0])

	if err == engine.ErrKeyNotFound {
		return fmt.Sprintf("record with key \"%s\" not found", query.Args[0])
	}

	if err != nil {
		d.logger.Error("failed to get value", zap.String("key", query.Args[0]), zap.Error(err))
		return fmt.Sprintf("error: %s", err.Error())
	}

	return val
}

func (d *Database) handleSetQuery(query *compute.Query) string {
	args := query.Args
	err := d.storage.Set(args[0], args[1])

	if err != nil {
		d.logger.Error("failed to set value", zap.String("key", args[0]), zap.String("value", args[1]), zap.Error(err))
		return fmt.Sprintf("error: %s", err.Error())
	}

	return "ok"
}

func (d *Database) handleDelQuery(query *compute.Query) string {
	err := d.storage.Del(query.Args[0])

	if err != nil {
		d.logger.Error("failed to delete value", zap.String("key", query.Args[0]), zap.Error(err))
		return fmt.Sprintf("error: %s", err.Error())
	}

	return "ok"
}
