package server

import (
	"net/http"
	"runtime"
	"time"

	"github.com/LarsArtmann/template-GoReleaser/internal/server/service"
	"github.com/LarsArtmann/template-GoReleaser/web/templates"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

// Health check handler
func (s *Server) healthHandler(c *gin.Context) {
	health := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"uptime":    time.Since(time.Now()).String(),
		"version":   "dev", // This would come from build info
	}

	c.JSON(http.StatusOK, health)
}

// Metrics handler
func (s *Server) metricsHandler(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := map[string]interface{}{
		"memory": map[string]interface{}{
			"alloc":       m.Alloc,
			"total_alloc": m.TotalAlloc,
			"sys":         m.Sys,
			"num_gc":      m.NumGC,
		},
		"goroutines": runtime.NumGoroutine(),
		"timestamp":  time.Now().UTC(),
	}

	c.JSON(http.StatusOK, metrics)
}

// Config validation handler
func (s *Server) validateConfigHandler(c *gin.Context) {
	var request struct {
		ConfigPath string `json:"config_path"`
		Content    string `json:"content,omitempty"`
		Config     string `form:"config"` // For form submissions
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Use content from form if available
	if request.Config != "" {
		request.Content = request.Config
	}

	// Get validation service
	validationService, err := do.InvokeAs[*service.ValidationService](s.injector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get validation service",
		})
		return
	}

	result := validationService.ValidateConfig(request.ConfigPath, request.Content)
	result.Timestamp = time.Now().UTC().Format(time.RFC3339)

	c.JSON(http.StatusOK, result)
}

// Environment validation handler
func (s *Server) validateEnvironmentHandler(c *gin.Context) {
	// Get validation service
	validationService, err := do.InvokeAs[*service.ValidationService](s.injector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get validation service",
		})
		return
	}

	result := validationService.ValidateEnvironment()
	result.Timestamp = time.Now().UTC().Format(time.RFC3339)

	c.JSON(http.StatusOK, result)
}

// Validation status handler
func (s *Server) validateStatusHandler(c *gin.Context) {
	// Get validation service
	validationService, err := do.InvokeAs[*service.ValidationService](s.injector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get validation service",
		})
		return
	}

	status := validationService.GetValidationStatus()
	status["last_check"] = time.Now().UTC().Format(time.RFC3339)

	c.JSON(http.StatusOK, status)
}

// Web UI handlers
func (s *Server) indexHandler(c *gin.Context) {
	// Get template service
	templateService, err := do.InvokeAs[*service.TemplateService](s.injector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get template service",
		})
		return
	}

	// Render template
	if err := templateService.RenderTempl(c, templates.Index()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to render template",
		})
	}
}

func (s *Server) configHandler(c *gin.Context) {
	// Get template service
	templateService, err := do.InvokeAs[*service.TemplateService](s.injector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get template service",
		})
		return
	}

	// Render template
	if err := templateService.RenderTempl(c, templates.Config()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to render template",
		})
	}
}

func (s *Server) updateConfigHandler(c *gin.Context) {
	var request struct {
		Config string `json:"config" form:"config"`
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Get config service
	configService, err := do.InvokeAs[*service.ConfigService](s.injector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get config service",
		})
		return
	}

	// Save config
	if err := configService.SaveConfigString("", request.Config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to save configuration",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Configuration updated successfully",
	})
}
