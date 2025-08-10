package validation

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserFriendlyError(t *testing.T) {
	originalErr := assert.AnError
	userErr := NewUserFriendlyError("technical message", "user message", originalErr)
	
	assert.Equal(t, "technical message", userErr.Message)
	assert.Equal(t, "user message", userErr.UserMessage)
	assert.Equal(t, originalErr, userErr.Cause)
	assert.Equal(t, "technical message", userErr.Error())
}

func TestValidationResult(t *testing.T) {
	t.Run("new validation result", func(t *testing.T) {
		result := NewValidationResult()
		
		assert.True(t, result.Valid)
		assert.Empty(t, result.Issues)
		assert.Empty(t, result.Warnings)
		assert.NotZero(t, result.Timestamp)
		assert.NotNil(t, result.Metadata)
	})
	
	t.Run("add error", func(t *testing.T) {
		result := NewValidationResult()
		result.AddError("test_field", "technical message", "user message")
		
		assert.False(t, result.Valid)
		assert.Len(t, result.Issues, 1)
		assert.Empty(t, result.Warnings)
		
		issue := result.Issues[0]
		assert.Equal(t, "test_field", issue.Field)
		assert.Equal(t, "technical message", issue.Message)
		assert.Equal(t, "user message", issue.UserMessage)
		assert.Equal(t, SeverityError, issue.Severity)
	})
	
	t.Run("add warning", func(t *testing.T) {
		result := NewValidationResult()
		result.AddWarning("test_field", "technical message", "user message")
		
		assert.True(t, result.Valid) // Warnings don't invalidate
		assert.Empty(t, result.Issues)
		assert.Len(t, result.Warnings, 1)
		
		warning := result.Warnings[0]
		assert.Equal(t, SeverityWarning, warning.Severity)
	})
	
	t.Run("add critical", func(t *testing.T) {
		result := NewValidationResult()
		result.AddCritical("test_field", "technical message", "user message")
		
		assert.False(t, result.Valid)
		assert.Len(t, result.Issues, 1)
		
		issue := result.Issues[0]
		assert.Equal(t, SeverityCritical, issue.Severity)
	})
	
	t.Run("has critical issues", func(t *testing.T) {
		result := NewValidationResult()
		assert.False(t, result.HasCriticalIssues())
		
		result.AddCritical("test", "msg", "user msg")
		assert.True(t, result.HasCriticalIssues())
	})
	
	t.Run("get issues by severity", func(t *testing.T) {
		result := NewValidationResult()
		result.AddError("field1", "msg1", "user1")
		result.AddWarning("field2", "msg2", "user2")
		result.AddCritical("field3", "msg3", "user3")
		
		criticalIssues := result.GetIssuesBySeverity(SeverityCritical)
		assert.Len(t, criticalIssues, 1)
		assert.Equal(t, "field3", criticalIssues[0].Field)
		
		warningIssues := result.GetIssuesBySeverity(SeverityWarning)
		assert.Len(t, warningIssues, 1)
		assert.Equal(t, "field2", warningIssues[0].Field)
	})
}

func TestEnvironmentValidationResult(t *testing.T) {
	t.Run("new environment validation result", func(t *testing.T) {
		result := NewEnvironmentValidationResult()
		
		assert.True(t, result.Valid)
		assert.Empty(t, result.CriticalMissing)
		assert.Empty(t, result.OptionalMissing)
		assert.NotNil(t, result.ValidatedVariables)
		assert.Equal(t, StatusReady, result.ValidationStatus)
	})
	
	t.Run("add missing critical", func(t *testing.T) {
		result := NewEnvironmentValidationResult()
		result.AddMissingCritical("GITHUB_TOKEN")
		
		assert.False(t, result.Valid)
		assert.Contains(t, result.CriticalMissing, "GITHUB_TOKEN")
		assert.Equal(t, StatusNeedsSetup, result.ValidationStatus)
	})
	
	t.Run("add missing optional", func(t *testing.T) {
		result := NewEnvironmentValidationResult()
		result.AddMissingOptional("DOCKER_TOKEN")
		
		assert.True(t, result.Valid) // Optional missing doesn't invalidate
		assert.Contains(t, result.OptionalMissing, "DOCKER_TOKEN")
		assert.Equal(t, StatusHasIssues, result.ValidationStatus)
	})
}

func TestConfigAnalysisResult(t *testing.T) {
	result := NewConfigAnalysisResult()
	
	assert.Empty(t, result.ConfigFiles)
	assert.Empty(t, result.ExtractedVariables)
	assert.Empty(t, result.MissingInExample)
	assert.Empty(t, result.UnusedInExample)
	assert.NotZero(t, result.Timestamp)
}

func TestValidationReport_JSON(t *testing.T) {
	// Test JSON marshaling of validation types
	t.Run("validation result JSON", func(t *testing.T) {
		result := NewValidationResult()
		result.AddError("test_field", "test message", "user message")
		
		jsonData, err := result.ToJSON()
		require.NoError(t, err)
		
		// Verify it's valid JSON
		var parsed map[string]interface{}
		err = json.Unmarshal(jsonData, &parsed)
		require.NoError(t, err)
		
		assert.False(t, parsed["valid"].(bool))
		assert.NotEmpty(t, parsed["issues"])
	})
	
	t.Run("environment validation result JSON", func(t *testing.T) {
		result := NewEnvironmentValidationResult()
		result.AddMissingCritical("GITHUB_TOKEN")
		
		jsonData, err := result.ToJSON()
		require.NoError(t, err)
		
		var parsed map[string]interface{}
		err = json.Unmarshal(jsonData, &parsed)
		require.NoError(t, err)
		
		assert.False(t, parsed["valid"].(bool))
		assert.Equal(t, "needs_setup", parsed["validation_status"])
	})
	
	t.Run("validation report JSON", func(t *testing.T) {
		report := &ValidationReport{
			Environment:    NewEnvironmentValidationResult(),
			ConfigAnalysis: NewConfigAnalysisResult(),
			OverallStatus:  StatusOK,
			Summary: ValidationSummary{
				TotalChecks: 10,
				Errors:      0,
				Warnings:    0,
			},
			Timestamp:          time.Now(),
			RecommendedActions: []string{"All good!"},
		}
		
		jsonData, err := report.ToJSON()
		require.NoError(t, err)
		
		var parsed map[string]interface{}
		err = json.Unmarshal(jsonData, &parsed)
		require.NoError(t, err)
		
		assert.Equal(t, "ok", parsed["overall_status"])
		assert.NotNil(t, parsed["environment"])
		assert.NotNil(t, parsed["config_analysis"])
	})
}

func TestValidationSeverityConstants(t *testing.T) {
	// Test that severity constants are defined correctly
	assert.Equal(t, ValidationSeverity("critical"), SeverityCritical)
	assert.Equal(t, ValidationSeverity("error"), SeverityError)
	assert.Equal(t, ValidationSeverity("warning"), SeverityWarning)
	assert.Equal(t, ValidationSeverity("info"), SeverityInfo)
}

func TestVariableCategoryConstants(t *testing.T) {
	// Test that category constants are defined correctly
	assert.Equal(t, VariableCategory("github"), CategoryGitHub)
	assert.Equal(t, VariableCategory("docker"), CategoryDocker)
	assert.Equal(t, VariableCategory("cloud"), CategoryCloud)
	assert.Equal(t, VariableCategory("signing"), CategorySigning)
	assert.Equal(t, VariableCategory("notification"), CategoryNotification)
	assert.Equal(t, VariableCategory("general"), CategoryGeneral)
	assert.Equal(t, VariableCategory("artifacts"), CategoryArtifacts)
}

func TestValidationStatusConstants(t *testing.T) {
	// Test that validation status constants are defined correctly
	assert.Equal(t, ValidationStatus("ok"), StatusOK)
	assert.Equal(t, ValidationStatus("warnings"), StatusWarnings)
	assert.Equal(t, ValidationStatus("errors"), StatusErrors)
	assert.Equal(t, ValidationStatus("critical_errors"), StatusCriticalErrors)
}

func TestEnvironmentValidationStatusConstants(t *testing.T) {
	// Test that environment validation status constants are defined correctly
	assert.Equal(t, EnvironmentValidationStatus("ready"), StatusReady)
	assert.Equal(t, EnvironmentValidationStatus("needs_setup"), StatusNeedsSetup)
	assert.Equal(t, EnvironmentValidationStatus("has_issues"), StatusHasIssues)
}

func TestVariableStatus(t *testing.T) {
	status := VariableStatus{
		Present: true,
		Valid:   false,
		Masked:  "gh***56",
		Issue:   "Invalid token format",
	}
	
	assert.True(t, status.Present)
	assert.False(t, status.Valid)
	assert.Equal(t, "gh***56", status.Masked)
	assert.Equal(t, "Invalid token format", status.Issue)
}