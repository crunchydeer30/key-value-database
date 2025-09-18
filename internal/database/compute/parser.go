package compute

import (
	"errors"
	"strings"

	"go.uber.org/zap"
)

type Parser struct {
	logger *zap.Logger
}

var ErrInvalidQuery = errors.New("invalid query")

const (
	maxArgs = 3
)

func NewParser(logger *zap.Logger) (*Parser, error) {
	if logger == nil {
		return nil, errors.New("no logger provided")
	}

	return &Parser{
		logger: logger,
	}, nil
}

func (p *Parser) Parse(queryStr string) (*Command, error) {
	queryStr = strings.TrimSpace(queryStr)
	parts := strings.Fields(queryStr)

	if len(parts) == 0 {
		p.logger.Debug("no tokens in query", zap.String("query", queryStr))
		return nil, ErrInvalidQuery
	}

	if len(parts) > maxArgs {
		p.logger.Debug("too many args in query", zap.String("query", queryStr))
		return nil, ErrInvalidNumberOfArgs
	}

	command := Command{
		Name: CommandName(parts[0]),
		Args: parts[1:],
	}

	err := command.validate()
	if err != nil {
		p.logger.Debug("invalid query", zap.String("query", queryStr), zap.Error(err))
		return nil, err
	}

	return &command, nil
}
