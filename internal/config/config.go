package config

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var (
	ErrValidationFailed = errors.New("validation failed")
	ErrReadConfigFailed = errors.New("read config failed")
	ErrUnmarshalFailed  = errors.New("unmarshal failed")
)

type Config struct {
	Engine EngineConfig `mapstructure:"engine"`
	Logger LoggerConfig `mapstructure:"logger"`
}

type EngineConfig struct {
	Type string `validate:"required,oneof=in_memory"`
}

type LoggerConfig struct {
	Level  string `validate:"required,oneof=debug info warn error"`
	Output string `validate:"required"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Join(ErrReadConfigFailed, err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Join(ErrUnmarshalFailed, err)
	}

	v := validator.New()
	if err := v.Struct(cfg); err != nil {
		return nil, errors.Join(ErrValidationFailed, err)
	}

	return &cfg, nil
}
