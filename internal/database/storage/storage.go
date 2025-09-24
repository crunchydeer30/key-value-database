package storage

type Engine interface {
	Set(key, value string) error
	Get(key, value string) (string, bool)
	Del(key string) error
}
