package service

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type TemplateService struct{}

func NewTemplateService(injector do.Injector) (*TemplateService, error) {
	return &TemplateService{}, nil
}

func (s *TemplateService) RenderTempl(c *gin.Context, component templ.Component) error {
	c.Header("Content-Type", "text/html")
	return component.Render(c.Request.Context(), c.Writer)
}

func (s *TemplateService) RenderTemplWithStatus(c *gin.Context, status int, component templ.Component) error {
	c.Status(status)
	c.Header("Content-Type", "text/html")
	return component.Render(c.Request.Context(), c.Writer)
}

// Helper function for rendering HTML templates (fallback)
func (s *TemplateService) RenderHTML(c *gin.Context, name string, data gin.H) {
	c.HTML(http.StatusOK, name, data)
}

// HTML template functions for dynamic content
func (s *TemplateService) GetTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatBytes": func(bytes int64) string {
			const unit = 1024
			if bytes < unit {
				return fmt.Sprintf("%d B", bytes)
			}
			div, exp := int64(unit), 0
			for n := bytes / unit; n >= unit; n /= unit {
				div *= unit
				exp++
			}
			return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
		},
	}
}