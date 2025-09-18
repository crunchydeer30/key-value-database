package compute

import (
	"errors"

	"go.uber.org/zap"
)

type Compute struct {
	Parser *Parser
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
		Parser: parser,
		logger: logger,
	}, nil
}
