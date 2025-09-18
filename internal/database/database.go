package database

import (
	"errors"
	"kv-db/internal/database/compute"

	"go.uber.org/zap"
)

type Database struct {
	compute *compute.Compute
	logger  *zap.Logger
}

func NewDatabase(logger *zap.Logger) (*Database, error) {
	if logger == nil {
		return nil, errors.New("no logger is provided")
	}

	compute, err := compute.NewCompute(logger)

	if err != nil {
		return nil, err
	}

	return &Database{
		compute: compute,
		logger:  logger,
	}, nil
}

func (d *Database) HandleQuery(queryStr string) error {
	d.logger.Debug("received query", zap.String("query", queryStr))

	_, err := d.compute.Parser.Parse(queryStr)
	if err != nil {
		return err
	}

	return nil
}
