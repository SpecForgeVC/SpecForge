package mcp

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/logger"
	"go.uber.org/zap"
)

// Server represents the MCP server instance
type Server struct {
	httpServer *http.Server
	handlers   *Handlers
	router     *Router
	mu         sync.Mutex
	running    bool
	config     Config
}

type Config struct {
	Port         int
	BindAddress  string
	Enabled      bool
	AuthRequired bool
	Token        string
	TokenService app.MCPTokenService
}

func NewServer(config Config, handlers *Handlers) *Server {
	s := &Server{
		config:   config,
		handlers: handlers,
	}
	s.router = NewRouter(s)
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.BindAddress, config.Port),
		Handler: s.router,
	}
	return s
}

// Start starts the MCP server in a background goroutine
func (s *Server) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("server already running")
	}
	s.running = true
	s.mu.Unlock()

	go func() {
		logger.MCPInfo("Starting MCP Server", zap.String("addr", s.httpServer.Addr))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.MCPError("MCP Server failed", zap.Error(err))
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
		}
	}()

	return nil
}

// Stop gracefully shuts down the MCP server
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	logger.MCPInfo("Stopping MCP Server")
	s.running = false
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

func (s *Server) GetConfig() Config {
	return s.config
}
