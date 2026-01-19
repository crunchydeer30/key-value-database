package network

import (
	"go.uber.org/zap"
)

type TCPServerOption func(*TCPServer)

func WithMaxConnections(max int) TCPServerOption {
	return func(s *TCPServer) {
		s.maxConnections = max
	}
}

func WithLogger(logger *zap.Logger) TCPServerOption {
	return func(s *TCPServer) {
		s.logger = logger
	}
}

func WithMaxMessageSize(max int) TCPServerOption {
	return func(s *TCPServer) {
		s.maxMessageSize = max
	}
}
