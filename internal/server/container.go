package server

import (
	"github.com/LarsArtmann/template-GoReleaser/internal/server/service"
	"github.com/samber/do/v2"
)

// Container holds all application dependencies
type Container struct {
	injector *do.RootScope
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	injector := do.New()
	container := &Container{injector: injector}

	// Register services
	container.registerServices()

	return container
}

// registerServices registers all application services
func (c *Container) registerServices() {
	// Register configuration service
	do.ProvideTransient(c.injector, service.NewConfigService)

	// Register validation service
	do.ProvideTransient(c.injector, service.NewValidationService)

	// Register template service
	do.ProvideTransient(c.injector, service.NewTemplateService)
}

// GetInjector returns the dependency injector
func (c *Container) GetInjector() *do.RootScope {
	return c.injector
}

// Shutdown cleans up the container
func (c *Container) Shutdown() error {
	return c.injector.Shutdown()
}
