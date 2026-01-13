package network

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/crunchydeer30/key-value-database/internal/sync"
	"go.uber.org/zap"
)

type TCPServer struct {
	listener       net.Listener
	sem            *sync.Semaphore
	maxConnections int
	logger         *zap.Logger
	handler        Handler
}

type Handler func([]byte) []byte

func NewTCPServer(addr string, handler Handler, opts ...TCPServerOption) (*TCPServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address %s: %w", addr, err)
	}

	//nolint:exhaustruct
	s := &TCPServer{
		listener: listener,
		handler:  handler,
		logger:   zap.NewNop(),
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.maxConnections > 0 {
		s.sem = sync.NewSemaphore(s.maxConnections)
	}

	return s, nil
}

func (s *TCPServer) Serve() {
	s.logger.Info("starting TCP server", zap.String("address", s.listener.Addr().String()))

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Error("failed to accept connection", zap.Error(err))
			continue
		}

		s.logger.Debug("accepted connection", zap.String("address", conn.RemoteAddr().String()))

		if s.sem != nil {
			s.sem.Acquire()
		}

		go func(conn net.Conn) {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Error("recovered from panic", zap.Any("panic", r))
				}
			}()

			defer func() {
				if err := conn.Close(); err != nil {
					s.logger.Error("failed to close connection", zap.Error(err))
				}
			}()

			if s.sem != nil {
				defer s.sem.Release()
			}

			err := conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
			if err != nil {
				s.logger.Error("failed to set read deadline", zap.Error(err))
				return
			}

			s.handle(conn)
		}(conn)
	}
}

func (s *TCPServer) handle(conn net.Conn) {
	r := bufio.NewReader(conn)

	for {
		data, err := r.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			s.logger.Debug("connection closed", zap.String("address", conn.RemoteAddr().String()))
			return
		}
		if err != nil {
			s.logger.Error("failed to read message", zap.Error(err))
			return
		}

		result := s.handler(data)
		result = append(result, '\n')

		if _, err := conn.Write(result); err != nil {
			s.logger.Error("failed to write response", zap.Error(err))
			return
		}
	}
}
