package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateGitHubToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid classic token",
			token:   "ghp_abcdefghijklmnopqrstuvwxyz1234567890",
			wantErr: false,
		},
		{
			name:    "valid fine-grained token",
			token:   "github_pat_11ABCDEFG0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			token:   "invalid_token",
			wantErr: true,
		},
		{
			name:    "too short classic",
			token:   "ghp_abc",
			wantErr: true,
		},
		{
			name:    "too long classic",
			token:   "ghp_abcdefghijklmnopqrstuvwxyz1234567890123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateGitHubToken(tt.token)
			if tt.wantErr {
				assert.True(t, result.IsError())
			} else {
				assert.True(t, result.IsOk())
			}
		})
	}
}

func TestValidateDockerToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid docker token",
			token:   "dckr_pat_abcdefghijklmnopqrstuvwxyz123456",
			wantErr: false,
		},
		{
			name:    "valid docker token with hyphens and underscores",
			token:   "dckr_pat_abc-def_ghi-jkl_mnopqrstuvwxyz123456",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid prefix",
			token:   "docker_pat_abcdefghijklmnopqrstuvwxyz123456",
			wantErr: true,
		},
		{
			name:    "too short",
			token:   "dckr_pat_abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateDockerToken(tt.token)
			if tt.wantErr {
				assert.True(t, result.IsError())
			} else {
				assert.True(t, result.IsOk())
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "valid email with plus",
			email:   "user+tag@example.com",
			wantErr: false,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "missing @",
			email:   "userexample.com",
			wantErr: true,
		},
		{
			name:    "missing domain",
			email:   "user@",
			wantErr: true,
		},
		{
			name:    "missing user",
			email:   "@example.com",
			wantErr: true,
		},
		{
			name:    "invalid domain",
			email:   "user@example",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateEmail(tt.email)
			if tt.wantErr {
				assert.True(t, result.IsError())
			} else {
				assert.True(t, result.IsOk())
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid https URL",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "valid http URL",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "valid URL with path",
			url:     "https://example.com/path/to/resource",
			wantErr: false,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "invalid scheme",
			url:     "ftp://example.com",
			wantErr: true,
		},
		{
			name:    "no scheme",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "no host",
			url:     "https://",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateURL(tt.url)
			if tt.wantErr {
				assert.True(t, result.IsError())
			} else {
				assert.True(t, result.IsOk())
			}
		})
	}
}

func TestValidateAWSBucketName(t *testing.T) {
	tests := []struct {
		name       string
		bucketName string
		wantErr    bool
	}{
		{
			name:       "valid bucket name",
			bucketName: "my-valid-bucket",
			wantErr:    false,
		},
		{
			name:       "valid bucket name with numbers",
			bucketName: "bucket123",
			wantErr:    false,
		},
		{
			name:       "empty bucket name",
			bucketName: "",
			wantErr:    true,
		},
		{
			name:       "too short",
			bucketName: "ab",
			wantErr:    true,
		},
		{
			name:       "too long",
			bucketName: "this-is-a-very-long-bucket-name-that-exceeds-the-limit-of-63-chars",
			wantErr:    true,
		},
		{
			name:       "uppercase letters",
			bucketName: "My-Bucket",
			wantErr:    true,
		},
		{
			name:       "consecutive hyphens",
			bucketName: "my--bucket",
			wantErr:    true,
		},
		{
			name:       "consecutive periods",
			bucketName: "my..bucket",
			wantErr:    true,
		},
		{
			name:       "starts with hyphen",
			bucketName: "-mybucket",
			wantErr:    true,
		},
		{
			name:       "ends with hyphen",
			bucketName: "mybucket-",
			wantErr:    true,
		},
		{
			name:       "ip address format",
			bucketName: "192.168.1.1",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateAWSBucketName(tt.bucketName)
			if tt.wantErr {
				assert.True(t, result.IsError())
			} else {
				assert.True(t, result.IsOk())
			}
		})
	}
}

func TestValidateGCSBucketName(t *testing.T) {
	tests := []struct {
		name       string
		bucketName string
		wantErr    bool
	}{
		{
			name:       "valid bucket name",
			bucketName: "my-valid-bucket",
			wantErr:    false,
		},
		{
			name:       "valid with underscores and dots",
			bucketName: "my_bucket.example",
			wantErr:    false,
		},
		{
			name:       "empty bucket name",
			bucketName: "",
			wantErr:    true,
		},
		{
			name:       "too short",
			bucketName: "ab",
			wantErr:    true,
		},
		{
			name:       "too long",
			bucketName: "this-is-a-very-long-bucket-name-that-exceeds-the-limit-of-63-chars",
			wantErr:    true,
		},
		{
			name:       "uppercase letters",
			bucketName: "My-Bucket",
			wantErr:    true,
		},
		{
			name:       "starts with hyphen",
			bucketName: "-mybucket",
			wantErr:    true,
		},
		{
			name:       "ends with hyphen",
			bucketName: "mybucket-",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateGCSBucketName(tt.bucketName)
			if tt.wantErr {
				assert.True(t, result.IsError())
			} else {
				assert.True(t, result.IsOk())
			}
		})
	}
}

func TestValidateAzureStorageName(t *testing.T) {
	tests := []struct {
		name        string
		storageName string
		wantErr     bool
	}{
		{
			name:        "valid storage name",
			storageName: "mystorage123",
			wantErr:     false,
		},
		{
			name:        "minimum length",
			storageName: "abc",
			wantErr:     false,
		},
		{
			name:        "maximum length",
			storageName: "abcdefghijklmnopqrstuvwx",
			wantErr:     false,
		},
		{
			name:        "empty storage name",
			storageName: "",
			wantErr:     true,
		},
		{
			name:        "too short",
			storageName: "ab",
			wantErr:     true,
		},
		{
			name:        "too long",
			storageName: "abcdefghijklmnopqrstuvwxy",
			wantErr:     true,
		},
		{
			name:        "uppercase letters",
			storageName: "MyStorage",
			wantErr:     true,
		},
		{
			name:        "with hyphens",
			storageName: "my-storage",
			wantErr:     true,
		},
		{
			name:        "with underscores",
			storageName: "my_storage",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateAzureStorageName(tt.storageName)
			if tt.wantErr {
				assert.True(t, result.IsError())
			} else {
				assert.True(t, result.IsOk())
			}
		})
	}
}

func TestValidateHostname(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		wantErr  bool
	}{
		{
			name:     "valid hostname",
			hostname: "example.com",
			wantErr:  false,
		},
		{
			name:     "valid subdomain",
			hostname: "api.example.com",
			wantErr:  false,
		},
		{
			name:     "valid with numbers",
			hostname: "server1.example.com",
			wantErr:  false,
		},
		{
			name:     "empty hostname",
			hostname: "",
			wantErr:  true,
		},
		{
			name:     "no TLD",
			hostname: "example",
			wantErr:  true,
		},
		{
			name:     "just TLD",
			hostname: ".com",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateHostname(tt.hostname)
			if tt.wantErr {
				assert.True(t, result.IsError())
			} else {
				assert.True(t, result.IsOk())
			}
		})
	}
}

func TestValidateFilePath(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "testfile.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "valid existing file",
			filePath: testFile,
			wantErr:  false,
		},
		{
			name:     "empty path",
			filePath: "",
			wantErr:  true,
		},
		{
			name:     "non-existent file",
			filePath: filepath.Join(tmpDir, "nonexistent.txt"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateFilePath(tt.filePath)
			if tt.wantErr {
				assert.True(t, result.IsError())
			} else {
				assert.True(t, result.IsOk())
			}
		})
	}
}

func TestValidateEnvironmentVariable(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		value    string
		envVar   EnvironmentVariable
		wantCode string
		wantSev  ValidationSeverity
	}{
		{
			name:    "missing required variable",
			varName: "GITHUB_TOKEN",
			value:   "",
			envVar: EnvironmentVariable{
				Name:     "GITHUB_TOKEN",
				Required: true,
			},
			wantCode: "MISSING_REQUIRED",
			wantSev:  SeverityCritical,
		},
		{
			name:    "missing optional variable",
			varName: "DOCKER_TOKEN",
			value:   "",
			envVar: EnvironmentVariable{
				Name:     "DOCKER_TOKEN",
				Required: false,
			},
			wantCode: "MISSING_OPTIONAL",
			wantSev:  SeverityWarning,
		},
		{
			name:    "placeholder value",
			varName: "GITHUB_TOKEN",
			value:   "your-token-here",
			envVar: EnvironmentVariable{
				Name:     "GITHUB_TOKEN",
				Required: true,
			},
			wantCode: "PLACEHOLDER_VALUE",
			wantSev:  SeverityWarning,
		},
		{
			name:    "valid value",
			varName: "GITHUB_TOKEN",
			value:   "ghp_validtokenformathere1234567890123456",
			envVar: EnvironmentVariable{
				Name:     "GITHUB_TOKEN",
				Required: true,
			},
			wantCode: "", // No issue
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issue := ValidateEnvironmentVariable(tt.varName, tt.value, tt.envVar)
			
			if tt.wantCode == "" {
				// Expecting no issue
				assert.True(t, IsValidationIssueEmpty(issue))
			} else {
				// Expecting an issue
				assert.False(t, IsValidationIssueEmpty(issue))
				assert.Equal(t, tt.wantCode, issue.Code)
				assert.Equal(t, tt.wantSev, issue.Severity)
			}
		})
	}
}

func TestMaskValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "short value",
			value:    "abc",
			expected: "***",
		},
		{
			name:     "normal value",
			value:    "secrettoken123",
			expected: "se***23",
		},
		{
			name:     "long value",
			value:    "verylongsecrettokenvalue",
			expected: "ve***ue",
		},
		{
			name:     "empty value",
			value:    "",
			expected: "***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskValue(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetValidators(t *testing.T) {
	validators := GetValidators()
	
	// Check that all expected validators are present
	expectedValidators := []string{
		"github_token",
		"docker_token",
		"email",
		"url",
		"aws_bucket_name",
		"gcs_bucket_name",
		"azure_storage_name",
		"hostname",
		"file_path",
	}
	
	for _, expected := range expectedValidators {
		assert.Contains(t, validators, expected, "Validator %s should be present", expected)
	}
	
	// Test that validators can be called
	for name, validator := range validators {
		t.Run("validator_"+name, func(t *testing.T) {
			// Just test that the validator can be called without panicking
			result := validator("test-value")
			// We don't care about the result, just that it doesn't panic
			assert.NotNil(t, result)
		})
	}
}