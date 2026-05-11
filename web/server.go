package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zeroibot/pack/clock"
	"github.com/zeroibot/pack/io"
)

type Config struct {
	Base     string
	Port     int
	CORSList []string
}

// LoadConfig loads the web Config from path
func LoadConfig(path string) (*Config, error) {
	cfg, err := io.ReadJSON[Config](path)
	if err != nil {
		return nil, err
	}
	if cfg.Base == "" || cfg.Port == 0 {
		return nil, fmt.Errorf("invalid web config")
	}
	// Make sure cfg.Base has no trailing slash
	cfg.Base = strings.TrimSuffix(cfg.Base, "/")
	return &cfg, nil
}

type Server struct {
	server *http.Server
}

// NewServer creates a new web Server
func NewServer(cfg *Config) *Server {
	server := new(http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
	})
	return new(Server{server})
}

// SetHandler sets the server's http handler
func (s *Server) SetHandler(handler http.Handler) {
	s.server.Handler = handler
}

// Run starts the web server
func (s *Server) Run() {
	message := fmt.Sprintf("Server started at %s", s.server.Addr)
	fmt.Printf("[INFO] (%s) %s\n", clock.DateTimeNow(), message)
	s.server.ListenAndServe()
}
