package engine

import "errors"

var ErrKeyNotFound = errors.New("key not found")

type Engine interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Del(key string) error
}
