package container

import (
	"github.com/LarsArtmann/template-GoReleaser/internal/services"
	"github.com/samber/do"
	"github.com/spf13/viper"
)

// Container wraps the samber/do injector
type Container struct {
	injector *do.Injector
}

// NewContainer creates a new DI container with all services registered
func NewContainer() *Container {
	injector := do.New()

	// Register core services
	registerCoreServices(injector)
	registerDomainServices(injector)

	return &Container{
		injector: injector,
	}
}

// GetInjector returns the underlying injector for service resolution
func (c *Container) GetInjector() *do.Injector {
	return c.injector
}

// Shutdown gracefully shuts down all services
func (c *Container) Shutdown() error {
	return c.injector.Shutdown()
}

// HealthCheck performs health checks on all registered services
func (c *Container) HealthCheck() map[string]error {
	return c.injector.HealthCheck()
}

// registerCoreServices registers fundamental application services
func registerCoreServices(injector *do.Injector) {
	// Configuration service
	do.Provide[*viper.Viper](injector, func(i *do.Injector) (*viper.Viper, error) {
		v := viper.New()
		v.SetConfigName(".goreleaser-cli")
		v.SetConfigType("yaml")
		v.AddConfigPath("$HOME")
		v.AddConfigPath(".")

		// Set defaults
		v.SetDefault("license.type", "MIT")
		v.SetDefault("cli.verbose", false)
		v.SetDefault("cli.colors", true)

		// Read config file if it exists
		_ = v.ReadInConfig() // Ignore error if config file doesn't exist

		return v, nil
	})
}

// registerDomainServices registers business logic services
func registerDomainServices(injector *do.Injector) {
	// Config service
	do.Provide[services.ConfigService](injector, func(i *do.Injector) (services.ConfigService, error) {
		viper := do.MustInvoke[*viper.Viper](i)
		return services.NewConfigService(viper), nil
	})

	// Validation service
	do.Provide[services.ValidationService](injector, func(i *do.Injector) (services.ValidationService, error) {
		configService := do.MustInvoke[services.ConfigService](i)
		return services.NewValidationService(configService), nil
	})

	// TODO: Add other domain services as they are implemented
	// - LicenseService
	// - VerificationService
}
