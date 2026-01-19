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
	Engine  EngineConfig  `mapstructure:"engine"`
	Logger  LoggerConfig  `mapstructure:"logger"`
	Network NetworkConfig `mapstructure:"network"`
}

type EngineConfig struct {
	Type string `validate:"required,oneof=in_memory"`
}

type LoggerConfig struct {
	Level  string `validate:"required,oneof=debug info warn error"`
	Output string `validate:"required"`
}

type NetworkConfig struct {
	Address        string `mapstructure:"address" validate:"required"`
	MaxConnections int    `mapstructure:"max_connections" validate:"omitempty,min=1"`
	MaxMessageSize int    `mapstructure:"max_message_size" validate:"min=1"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)

	viper.AutomaticEnv()
	viper.SetDefault("engine.type", "in_memory")
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("logger.output", "stdout")
	viper.SetDefault("network.address", "127.0.0.1:3223")
	viper.SetDefault("network.max_connections", 100)
	viper.SetDefault("network.max_message_size", 4096)

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
