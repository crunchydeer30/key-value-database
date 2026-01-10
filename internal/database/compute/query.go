package compute

import (
	"errors"
)

var (
	ErrInvalidNumberOfArgs = errors.New("invalid number of args")
	ErrUnknownCommand      = errors.New("unknown command")
)

type CommandName string

const (
	GET CommandName = "GET"
	SET CommandName = "SET"
	DEL CommandName = "DEL"
)

const (
	getCommandArgsCount = 1
	setCommandArgsCount = 2
	delCommandArgsCount = 1
)

type Query struct {
	Command CommandName
	Args    []string
}

func NewQuery(command CommandName, args []string) Query {
	return Query{
		Command: command,
		Args:    args,
	}
}

func (q *Query) validate() error {
	switch q.Command {
	case GET:
		if len(q.Args) != getCommandArgsCount {
			return ErrInvalidNumberOfArgs
		}
	case SET:
		if len(q.Args) != setCommandArgsCount {
			return ErrInvalidNumberOfArgs
		}
	case DEL:
		if len(q.Args) != delCommandArgsCount {
			return ErrInvalidNumberOfArgs
		}
	default:
		return ErrUnknownCommand
	}

	return nil
}
