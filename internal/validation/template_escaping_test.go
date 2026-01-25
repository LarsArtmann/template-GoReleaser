package validation

import (
	"testing"
)

// Helper type for escape function tests
type escapeTestCase struct {
	name     string
	input    string
	expected string
}

// Helper type for validation function tests
type validationTestCase struct {
	name     string
	input    string
	expected bool
}

// escapeFunc represents a function that takes a string and returns an escaped string
type escapeFunc func(string) string

// validationFunc represents a function that takes a string and returns a bool
type validationFunc func(string) bool

// testFunc represents a generic test function
type testFunc[T any] func(string) T

// valueTestCase represents a generic test case for value-returning functions
type valueTestCase[T any] struct {
	name     string
	input    string
	expected T
}

// runTests is a generic helper function to run table-driven tests
func runTests[T comparable](t *testing.T, funcName string, fn testFunc[T], tests []valueTestCase[T]) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fn(tt.input)
			if result != tt.expected {
				t.Errorf("%s() = %v, want %v", funcName, result, tt.expected)
			}
		})
	}
}

// runFuzzTest is a helper function to run fuzz tests for escape functions
func runFuzzTest(f *testing.F, seed []string, escaper escapeFunc) {
	for _, s := range seed {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		result := escaper(input)
		// Should not panic
		_ = result
	})
}

// Test case constants for validation functions
var (
	shellInjectionTestCases = []testCase[bool]{
		{"echo hello", "echo hello", false},
		{"rm -rf /", "rm -rf /", true},
		{"cat file | grep pattern", "cat file | grep pattern", true},
		{"command && rm file", "command && rm file", true},
		{"script.sh", "script.sh", false},
		{"$(rm file)", "$(rm file)", true},
		{"`rm file`", "`rm file`", true},
	}

	dockerLabelTestCases = []testCase[bool]{
		{"label", "label", true},
		{"my-label.v1", "my-label.v1", true},
		{"my_label", "my_label", true},
		{"empty", "", false},
		{"invalid@label", "invalid@label", false},
		{"label with spaces", "label with spaces", false},
		{"label/with/slashes", "label/with/slashes", false},
	}
)

// runBenchmark is a helper function to run benchmark tests for escape functions
func runBenchmark(b *testing.B, input string, escaper escapeFunc) {
	for b.Loop() {
		_ = escaper(input)
	}
}

// runEscapeTest is a helper function to run a complete escape test with a fresh escaper
func runEscapeTest(t *testing.T, funcName string, escaperMethod func(te *TemplateEscaper, input string) string, tests []testCase[string]) {
	te := NewTemplateEscaper()
	escaper := func(input string) string {
		return escaperMethod(te, input)
	}
	runTests(t, funcName, escaper, tests)
}

func TestTemplateEscaper_EscapeYAML(t *testing.T) {
	tests := []testCase[string]{
		{"Simple text", "hello", "hello"},
		{"Empty string", "", ""},
		{"String with colon", "name: value", "'name: value'"},
		{"String with space", " leading space", "leading space"},
		{"Number", "123", "'123'"},
		{"Boolean-like", "true", "true"},
		{"String starting with special", "!important", "'!important'"},
		{"Multi-line", "line1\nline2", "|-\nline1\n  line2"},
		{"Complex multi-line", "line1: value\nline2: value", "|-\nline1: value\n  line2: value"},
	}

	runEscapeTest(t, "EscapeYAML", (*TemplateEscaper).EscapeYAML, tests)
}

func TestTemplateEscaper_EscapeShell(t *testing.T) {
	tests := []testCase[string]{
		{"Simple text", "hello", "'hello'"},
		{"Empty string", "", ""},
		{"Single quotes", "don't panic", "'don''t panic'"},
		{"Safe characters", "my-app_v1.0", "'my-app_v1.0'"},
		{"Dangerous content", "rm -rf /", ""},     // Should be filtered
		{"Script injection", "; echo hacked", ""}, // Should be filtered
	}

	runEscapeTest(t, "EscapeShell", (*TemplateEscaper).EscapeShell, tests)
}

func TestTemplateEscaper_EscapeGitHubActions(t *testing.T) {
	tests := []testCase[string]{
		{"Simple text", "hello", "hello"},
		{"Expression syntax", "${{ github.repository }}", "'${{ '' }}${{ github.repository }}'"},
		{"Empty string", "", ""},
		{"Complex YAML", "name: value", "'name: value'"},
	}

	runEscapeTest(t, "EscapeGitHubActions", (*TemplateEscaper).EscapeGitHubActions, tests)
}

func TestTemplateEscaper_EscapeJSON(t *testing.T) {
	tests := []testCase[string]{
		{"Simple text", "hello", `"hello"`},
		{"Empty string", "", `""`},
		{"Quotes", "say \"hello\"", `"say \"hello\""`},
		{"Backslash", "path\\to\\file", `"path\\to\\file"`},
		{"Newline", "line1\nline2", `"line1\nline2"`},
		{"Tab", "col1\tcol2", `"col1\tcol2"`},
	}

	runEscapeTest(t, "EscapeJSON", (*TemplateEscaper).EscapeJSON, tests)
}

func TestTemplateEscaper_EscapeDockerLabel(t *testing.T) {
	tests := []testCase[string]{
		{"Simple text", "hello", "hello"},
		{"Empty string", "", ""},
		{"Valid characters", "my-app_v1.0", "my-app_v1.0"},
		{"Invalid characters", "my@app$", "my-app-"},
		{"Starts with number", "123label", "label-123label"},
		{"Starts with dot", ".hidden", "label-.hidden"},
		{"Starts with dash", "-dash", "label--dash"},
	}

	runEscapeTest(t, "EscapeDockerLabel", (*TemplateEscaper).EscapeDockerLabel, tests)
}

func TestTemplateEscaper_ValidateTemplateContent(t *testing.T) {
	te := NewTemplateEscaper()

	tests := []struct {
		name         string
		content      string
		templateType string
		wantErr      bool
	}{
		{"Safe YAML", "name: myapp", "yaml", false},
		{"YAML injection", "name: ${SCRIPT}", "yaml", false}, // Pattern doesn't include full variable
		{"Safe shell", "echo 'hello'", "shell", false},
		{"Shell injection", "rm -rf /", "shell", true},
		{"Safe GitHub Actions", "name: build", "github-actions", false},
		{"GitHub Actions injection", "${{ github.token }}", "github-actions", true},
		{"Script tag", "<script>alert('xss')</script", "any", true},
		{"JavaScript", "javascript:void(0)", "any", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := te.ValidateTemplateContent(tt.content, tt.templateType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTemplateContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLooksLikeNumber(t *testing.T) {
	tests := []testCase[bool]{
		{"123", "123", true},
		{"123.45", "123.45", true},
		{"-123", "-123", true},
		{"+123", "+123", true},
		{"1e10", "1e10", true},
		{"123abc", "123abc", false},
		{"abc123", "abc123", false},
		{"empty", "", false},
		{"-", "-", false},
		{".", ".", false},
	}

	runTests(t, "looksLikeNumber", looksLikeNumber, tests)
}

func TestContainsShellInjection(t *testing.T) {
	runTests(t, "containsShellInjection", containsShellInjection, shellInjectionTestCases)
}

func TestIsValidDockerLabel(t *testing.T) {
	runTests(t, "isValidDockerLabel", isValidDockerLabel, dockerLabelTestCases)
}

// Fuzzing tests for escaping functions.
func FuzzEscapeYAML(f *testing.F) {
	seed := []string{"hello", "name: value", "don't", "multiline\nstring"}
	te := NewTemplateEscaper()
	runFuzzTest(f, seed, te.EscapeYAML)
}

func FuzzEscapeShell(f *testing.F) {
	seed := []string{"hello", "rm -rf", "; echo", "command$(rm)"}
	te := NewTemplateEscaper()
	runFuzzTest(f, seed, te.EscapeShell)
}

func FuzzEscapeJSON(f *testing.F) {
	seed := []string{"hello", "quote's", `back\slash`, "newline\n"}
	te := NewTemplateEscaper()
	runFuzzTest(f, seed, te.EscapeJSON)
}

// Benchmark tests.
func BenchmarkEscapeYAML(b *testing.B) {
	te := NewTemplateEscaper()
	input := "my-project-name: value with 'quotes'"
	runBenchmark(b, input, te.EscapeYAML)
}

func BenchmarkEscapeShell(b *testing.B) {
	te := NewTemplateEscaper()
	input := "my-app with single 'quotes'"
	runBenchmark(b, input, te.EscapeShell)
}

func BenchmarkEscapeJSON(b *testing.B) {
	te := NewTemplateEscaper()
	input := `string with "quotes" and \backslashes`
	runBenchmark(b, input, te.EscapeJSON)
}
