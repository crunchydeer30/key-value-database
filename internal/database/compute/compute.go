package compute

import (
	"errors"

	"go.uber.org/zap"
)

type Compute struct {
	parser *Parser
	logger *zap.Logger
}

func NewCompute(logger *zap.Logger) (*Compute, error) {
	if logger == nil {
		return nil, errors.New("no logger provided")
	}

	parser, err := NewParser(logger)
	if err != nil {
		return nil, errors.Join(errors.New("failed to initialize parsers"), err)
	}
	logger.Debug("parser initialized")

	return &Compute{
		parser: parser,
		logger: logger,
	}, nil
}

func (c *Compute) Parse(queryStr string) (*Query, error) {
	return c.parser.Parse(queryStr)
}
