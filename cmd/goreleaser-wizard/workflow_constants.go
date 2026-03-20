package main

import "time"

const (
	defaultTimeout    = 30 * time.Minute
	fullWizardTimeout = 10 * time.Minute
	configOnlyTimeout = 5 * time.Minute
	validationTimeout = 2 * time.Minute
	migrationTimeout  = 15 * time.Minute
	updateTimeout     = 10 * time.Minute
)
