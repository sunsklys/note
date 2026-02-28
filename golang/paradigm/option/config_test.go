package option

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"
)

type Config struct {
	Protocol string
	Timeout  time.Duration
	Maxconns int
	TLS      *tls.Config
}

type Server struct {
	Addr string
	Port int
	Conf *Config
}

func NewServerConfig(addr string, port int, conf *Config) (*Server, error) {
	return &Server{
		Addr: addr,
		Port: port,
		Conf: conf,
	}, nil
}

func TestConfig(t *testing.T) {

	srv1, _ := NewServerConfig("localhost", 9000, nil)

	conf := Config{Protocol: "tcp", Timeout: 60 * time.Second}
	srv2, _ := NewServerConfig("locahost", 9000, &conf)

	fmt.Printf("%+v\n", srv1)
	fmt.Printf("%+v\n", srv2)
}
