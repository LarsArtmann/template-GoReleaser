package fixtures

// TestEnvironmentVars provides test environment variables for different scenarios
var TestEnvironmentVars = map[string]map[string]string{
	"minimal": {
		"GITHUB_TOKEN":    "test-token",
		"DOCKER_USERNAME": "testuser",
		"DOCKER_PASSWORD": "testpass",
	},
	"complete": {
		"GITHUB_TOKEN":                "test-token",
		"DOCKER_USERNAME":             "testuser",
		"DOCKER_PASSWORD":             "testpass",
		"GORELEASER_KEY":              "test-pro-key",
		"COSIGN_PRIVATE_KEY":          "test-cosign-key",
		"FURY_TOKEN":                  "test-fury-token",
		"CHOCOLATEY_API_KEY":          "test-choco-key",
		"AUR_KEY":                     "test-aur-key",
		"HOMEBREW_TAP_GITHUB_TOKEN":   "test-homebrew-token",
		"DISCORD_WEBHOOK_ID":          "123456789",
		"DISCORD_WEBHOOK_TOKEN":       "test-discord-token",
		"SLACK_WEBHOOK":               "https://hooks.slack.com/test",
		"TWITTER_CONSUMER_KEY":        "test-twitter-key",
		"TWITTER_CONSUMER_SECRET":     "test-twitter-secret",
		"TWITTER_ACCESS_TOKEN":        "test-twitter-token",
		"TWITTER_ACCESS_TOKEN_SECRET": "test-twitter-token-secret",
		"MASTODON_CLIENT_ID":          "test-mastodon-id",
		"MASTODON_CLIENT_SECRET":      "test-mastodon-secret",
		"MASTODON_ACCESS_TOKEN":       "test-mastodon-token",
		"MASTODON_SERVER":             "https://mastodon.social",
		"REDDIT_APPLICATION_ID":       "test-reddit-id",
		"REDDIT_USERNAME":             "testuser",
		"REDDIT_PASSWORD":             "testpass",
		"TELEGRAM_TOKEN":              "123:test-telegram-token",
		"TELEGRAM_CHAT_ID":            "test-chat-id",
		"LINKEDIN_ACCESS_TOKEN":       "test-linkedin-token",
		"PROJECT_NAME":                "test-project",
		"PROJECT_DESCRIPTION":         "Test project for integration testing",
		"PROJECT_AUTHOR":              "Test Author",
		"PROJECT_EMAIL":               "test@example.com",
		"PROJECT_URL":                 "https://github.com/testuser/test-project",
		"LICENSE_TYPE":                "MIT",
	},
}

// ReadmeConfigs provides test readme configurations
var ReadmeConfigs = map[string]string{
	"minimal": `project:
  name: test-project
  description: Test project for integration testing

author:
  name: Test Author
  email: test@example.com

license:
  type: MIT

repository:
  url: https://github.com/testuser/test-project
`,
	"complete": `project:
  name: test-project
  description: Test project for integration testing
  version: 1.0.0

author:
  name: Test Author
  email: test@example.com
  url: https://example.com

license:
  type: MIT
  year: 2024

repository:
  url: https://github.com/testuser/test-project
  issues: https://github.com/testuser/test-project/issues
  wiki: https://github.com/testuser/test-project/wiki

features:
  - Cross-platform binaries
  - Docker images
  - Package manager distributions
  - Automatic releases

technologies:
  - Go
  - Docker
  - GitHub Actions
`,
}

// ExpectedLicenseContent provides expected license content for validation
var ExpectedLicenseContent = map[string]string{
	"MIT": `MIT License

Copyright (c) 2024 Test Author

Permission is hereby granted, free of charge, to any person obtaining a copy`,
	"Apache-2.0": `Apache License
Version 2.0, January 2004
http://www.apache.org/licenses/

TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION`,
	"BSD-3-Clause": `BSD 3-Clause License

Copyright (c) 2024, Test Author
All rights reserved.

Redistribution and use in source and binary forms`,
}

// GoModContent provides test go.mod content
const GoModContent = `module github.com/testuser/test-project

go 1.23

require github.com/stretchr/testify v1.10.0
`

// MainGoContent provides test main.go content
const MainGoContent = `package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	builtBy   = "unknown"
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
	
	fmt.Println("Hello from test-project!")
	fmt.Printf("Version: %s\n", version)
}

func printVersion() {
	fmt.Printf("Version:      %s\n", version)
	fmt.Printf("Commit:       %s\n", commit)
	fmt.Printf("Built:        %s\n", date)
	fmt.Printf("Built by:     %s\n", builtBy)
	fmt.Printf("Go version:   %s\n", runtime.Version())
	fmt.Printf("OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
`
