package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type Server struct {
	injector *do.RootScope
	router   *gin.Engine
	addr     string
}

type Config struct {
	Port     int
	Host     string
	Debug    bool
	StaticDir string
}

func New(injector *do.RootScope, config Config) *Server {
	// Set gin mode based on debug setting
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	
	// Add middleware
	router.Use(gin.Recovery())
	if config.Debug {
		router.Use(gin.Logger())
	}

	server := &Server{
		injector: injector,
		router:   router,
		addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.healthHandler)
	
	// Metrics endpoint
	s.router.GET("/metrics", s.metricsHandler)
	
	// Validation endpoints
	v1 := s.router.Group("/api/v1")
	{
		v1.POST("/validate/config", s.validateConfigHandler)
		v1.POST("/validate/environment", s.validateEnvironmentHandler)
		v1.GET("/validate/status", s.validateStatusHandler)
	}

	// Web UI routes
	s.router.GET("/", s.indexHandler)
	s.router.GET("/config", s.configHandler)
	s.router.POST("/config", s.updateConfigHandler)
	
	// Static files
	s.router.Static("/static", "./web/static")
}

func (s *Server) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", s.addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}

func (s *Server) Router() *gin.Engine {
	return s.router
}