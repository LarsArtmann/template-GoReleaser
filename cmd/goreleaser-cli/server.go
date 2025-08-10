package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/LarsArtmann/template-GoReleaser/internal/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the web server for GoReleaser configuration management",
	Long: `Start a web server that provides a modern web interface for managing
and validating GoReleaser configurations. The server includes:

- Interactive configuration editor with HTMX
- Real-time validation and syntax checking
- Health monitoring and metrics
- Environment variable validation
- API endpoints for automation

The web interface provides an intuitive way to create, edit, and validate
your GoReleaser configurations without needing to use the command line.`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Server configuration flags
	serverCmd.Flags().IntP("port", "p", 8080, "Port to run the server on")
	serverCmd.Flags().String("host", "localhost", "Host to bind the server to")
	serverCmd.Flags().Bool("debug", false, "Enable debug mode with verbose logging")
	serverCmd.Flags().String("static-dir", "./web/static", "Directory for static files")
}

func runServer(cmd *cobra.Command, args []string) {
	// Get configuration from flags
	port, _ := cmd.Flags().GetInt("port")
	host, _ := cmd.Flags().GetString("host")
	debug, _ := cmd.Flags().GetBool("debug")
	staticDir, _ := cmd.Flags().GetString("static-dir")

	fmt.Printf("🚀 Starting GoReleaser Configuration Server\n")
	fmt.Printf("   Host: %s\n", host)
	fmt.Printf("   Port: %d\n", port)
	fmt.Printf("   Debug: %t\n", debug)
	fmt.Printf("   Static Dir: %s\n", staticDir)
	fmt.Println()

	// Create dependency injection container
	container := server.NewContainer()
	defer func() {
		if err := container.Shutdown(); err != nil {
			fmt.Printf("Error shutting down container: %v\n", err)
		}
	}()

	// Create server configuration
	config := server.Config{
		Port:      port,
		Host:      host,
		Debug:     debug,
		StaticDir: staticDir,
	}

	// Create server instance
	srv := server.New(container.GetInjector(), config)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		
		fmt.Println("\n🛑 Shutdown signal received")
		cancel()
	}()

	// Start server
	fmt.Printf("🌐 Server starting at http://%s:%d\n", host, port)
	fmt.Printf("📋 Configuration Editor: http://%s:%d/config\n", host, port)
	fmt.Printf("🏥 Health Check: http://%s:%d/health\n", host, port)
	fmt.Printf("📊 Metrics: http://%s:%d/metrics\n", host, port)
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop the server")

	// Start the server (blocking)
	if err := srv.Start(ctx); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Server stopped gracefully")
}