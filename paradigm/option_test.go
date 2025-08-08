package paradigm

import (
	"testing"
	"time"
)

func TestOption(t *testing.T) {
	options := []Option{Protocol("udp"), Timeout(60 * time.Second), MaxConn(80)}
	srv := NewServer("127.0.0.1", 3306, options...)
	t.Logf("%+v\n", srv)
}

type Server struct {
	Addr string
	Port int
	Config
}

type Config struct {
	Protocol string
	Timeout  time.Duration
	MaxConn  int
}

func NewServer(addr string, port int, options ...Option) *Server {
	srv := &Server{
		Addr: addr,
		Port: port,
		Config: Config{
			Protocol: "tcp",
			Timeout:  30 * time.Second,
			MaxConn:  10,
		},
	}

	for _, option := range options {
		option(srv)
	}

	return srv
}

type Option func(*Server)

func Protocol(p string) Option {
	return func(s *Server) {
		s.Protocol = p
	}
}

func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Timeout = timeout
	}
}

func MaxConn(maxConn int) Option {
	return func(s *Server) {
		s.MaxConn = maxConn
	}
}
