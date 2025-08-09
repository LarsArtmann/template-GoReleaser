package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	version      = "dev"
	commit       = "none"
	date         = "unknown"
	builtBy      = "unknown"
	gitDescription = ""
	gitState     = ""
)

func main() {
	var (
		versionFlag = flag.Bool("version", false, "Print version information")
		healthFlag  = flag.Bool("health", false, "Health check")
	)
	
	flag.Parse()
	
	if *versionFlag {
		printVersion()
		os.Exit(0)
	}
	
	if *healthFlag {
		fmt.Println("OK")
		os.Exit(0)
	}
	
	fmt.Println("Hello from GoReleaser Template!")
	fmt.Printf("Version: %s\n", version)
}

func printVersion() {
	fmt.Printf("Version:      %s\n", version)
	fmt.Printf("Commit:       %s\n", commit)
	fmt.Printf("Built:        %s\n", date)
	fmt.Printf("Built by:     %s\n", builtBy)
	fmt.Printf("Go version:   %s\n", runtime.Version())
	fmt.Printf("OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
	
	if gitDescription != "" {
		fmt.Printf("Git describe: %s\n", gitDescription)
	}
	
	if gitState != "" {
		fmt.Printf("Git state:    %s\n", gitState)
	}
}