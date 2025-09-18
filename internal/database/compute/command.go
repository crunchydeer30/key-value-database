package compute

import (
	"errors"
)

var ErrInvalidNumberOfArgs = errors.New("invalid number of args")
var ErrUnknownCommand = errors.New("unknown command")

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

type Command struct {
	Name CommandName
	Args []string
}

func (c *Command) validate() error {
	switch c.Name {
	case GET:
		if len(c.Args) != getCommandArgsCount {
			return ErrInvalidNumberOfArgs
		}
	case SET:
		if len(c.Args) != setCommandArgsCount {
			return ErrInvalidNumberOfArgs
		}
	case DEL:
		if len(c.Args) != delCommandArgsCount {
			return ErrInvalidNumberOfArgs
		}
	default:
		return ErrUnknownCommand
	}

	return nil
}
