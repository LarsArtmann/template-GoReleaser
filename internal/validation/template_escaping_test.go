package validation

import (
	"testing"
)

// escapeFunc represents a function that takes a string and returns an escaped string.
type escapeFunc func(string) string

// testFunc represents a generic test function.
type testFunc[T any] func(string) T

// valueTestCase represents a generic test case for value-returning functions.
type valueTestCase[T any] struct {
	name     string
	input    string
	expected T
}

// runTests is a generic helper function to run table-driven tests.
func runTests[T comparable](
	t *testing.T,
	funcName string,
	fn testFunc[T],
	tests []valueTestCase[T],
) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fn(tt.input)
			if result != tt.expected {
				t.Errorf("%s() = %v, want %v", funcName, result, tt.expected)
			}
		})
	}
}

// runFuzzTest is a helper function to run fuzz tests for escape functions.
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

// stringTestCaseInput is the input type for String escape test case builders.
type stringTestCaseInput struct {
	Name     string
	Input    string
	Expected string
}

// boolTestCaseInput is the input type for boolean validation test case builders.
type boolTestCaseInput struct {
	Name     string
	Input    string
	Expected bool
}

// stringEscaperCase creates a stringTestCaseInput for escape function tests.
func stringEscaperCase(name, input, expected string) stringTestCaseInput {
	return stringTestCaseInput{Name: name, Input: input, Expected: expected}
}

// boolValidationCase creates a boolTestCaseInput for validation function tests.
func boolValidationCase(name, input string, expected bool) boolTestCaseInput {
	return boolTestCaseInput{Name: name, Input: input, Expected: expected}
}

// escapeStringTestCaseBuilder builds String escape test cases from variadic stringTestCaseInput.
func escapeStringTestCaseBuilder(cases ...stringTestCaseInput) []valueTestCase[string] {
	result := make([]valueTestCase[string], len(cases))
	for i, c := range cases {
		result[i] = valueTestCase[string]{name: c.Name, input: c.Input, expected: c.Expected}
	}

	return result
}

// escapeBoolTestCaseBuilder builds boolean validation test cases from variadic boolTestCaseInput.
func escapeBoolTestCaseBuilder(cases ...boolTestCaseInput) []valueTestCase[bool] {
	result := make([]valueTestCase[bool], len(cases))
	for i, c := range cases {
		result[i] = valueTestCase[bool]{name: c.Name, input: c.Input, expected: c.Expected}
	}

	return result
}

// Test case constants for escape functions.
var (
	shellEscapeTestCases = escapeStringTestCaseBuilder(
		stringEscaperCase("Simple text", "hello", "'hello'"),
		stringEscaperCase("Empty string", "", ""),
		stringEscaperCase("Single quotes", "don't panic", "'don''t panic'"),
		stringEscaperCase("Safe characters", "my-app_v1.0", "'my-app_v1.0'"),
		stringEscaperCase("Dangerous content", "rm -rf /", ""),
		stringEscaperCase("Script injection", "; echo hacked", ""),
	)

	jsonEscapeTestCases = escapeStringTestCaseBuilder(
		stringEscaperCase("Simple text", "hello", `"hello"`),
		stringEscaperCase("Empty string", "", `""`),
		stringEscaperCase("Quotes", "say \"hello\"", `"say \"hello\""`),
		stringEscaperCase("Backslash", "path\\to\\file", `"path\\to\\file"`),
		stringEscaperCase("Newline", "line1\nline2", `"line1\nline2"`),
		stringEscaperCase("Tab", "col1\tcol2", `"col1\tcol2"`),
	)

	yamlEscapeTestCases = escapeStringTestCaseBuilder(
		stringEscaperCase("Simple text", "hello", "hello"),
		stringEscaperCase("Empty string", "", ""),
		stringEscaperCase("String with colon", "name: value", "'name: value'"),
		stringEscaperCase("String with space", " leading space", "leading space"),
		stringEscaperCase("Number", "123", "'123'"),
		stringEscaperCase("Boolean-like", "true", "true"),
		stringEscaperCase("String starting with special", "!important", "'!important'"),
		stringEscaperCase("Multi-line", "line1\nline2", "|-\nline1\n  line2"),
		stringEscaperCase(
			"Complex multi-line",
			"line1: value\nline2: value",
			"|-\nline1: value\n  line2: value",
		),
	)

	githubActionsEscapeTestCases = escapeStringTestCaseBuilder(
		stringEscaperCase("Simple text", "hello", "hello"),
		stringEscaperCase(
			"Expression syntax",
			"${{ github.repository }}",
			"'${{ '' }}${{ github.repository }}'",
		),
		stringEscaperCase("Empty string", "", ""),
		stringEscaperCase("Complex YAML", "name: value", "'name: value'"),
	)

	dockerLabelEscapeTestCases = escapeStringTestCaseBuilder(
		stringEscaperCase("Simple text", "hello", "hello"),
		stringEscaperCase("Empty string", "", ""),
		stringEscaperCase("Valid characters", "my-app_v1.0", "my-app_v1.0"),
		stringEscaperCase("Invalid characters", "my@app$", "my-app-"),
		stringEscaperCase("Starts with number", "123label", "label-123label"),
		stringEscaperCase("Starts with dot", ".hidden", "label-.hidden"),
		stringEscaperCase("Starts with dash", "-dash", "label--dash"),
	)

	shellInjectionTestCases = escapeBoolTestCaseBuilder(
		boolValidationCase("echo hello", "echo hello", false),
		boolValidationCase("rm -rf /", "rm -rf /", true),
		boolValidationCase("cat file | grep pattern", "cat file | grep pattern", true),
		boolValidationCase("command && rm file", "command && rm file", true),
		boolValidationCase("script.sh", "script.sh", false),
		boolValidationCase("$(rm file)", "$(rm file)", true),
		boolValidationCase("`rm file`", "`rm file`", true),
	)

	dockerLabelTestCases = escapeBoolTestCaseBuilder(
		boolValidationCase("label", "label", true),
		boolValidationCase("my-label.v1", "my-label.v1", true),
		boolValidationCase("my_label", "my_label", true),
		boolValidationCase("empty", "", false),
		boolValidationCase("invalid@label", "invalid@label", false),
		boolValidationCase("label with spaces", "label with spaces", false),
		boolValidationCase("label/with/slashes", "label/with/slashes", false),
	)

	// LooksLikeNumberTestCases shared across number detection tests.
	looksLikeNumberTestCases = []valueTestCase[bool]{
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
)

// runBenchmark is a helper function to run benchmark tests for escape functions.
func runBenchmark(b *testing.B, input string, escaper escapeFunc) {
	for b.Loop() {
		_ = escaper(input)
	}
}

// runEscapeTest is a helper function to run a complete escape test with a fresh escaper.
func runEscapeTest(
	t *testing.T,
	funcName string,
	escaperMethod func(te *TemplateEscaper, input string) string,
	tests []valueTestCase[string],
) {
	te := NewTemplateEscaper()
	escaper := func(input string) string {
		return escaperMethod(te, input)
	}
	runTests(t, funcName, escaper, tests)
}

func TestTemplateEscaper_EscapeYAML(t *testing.T) {
	runEscapeTest(t, "EscapeYAML", (*TemplateEscaper).EscapeYAML, yamlEscapeTestCases)
}

func TestTemplateEscaper_EscapeShell(t *testing.T) {
	runEscapeTest(t, "EscapeShell", (*TemplateEscaper).EscapeShell, shellEscapeTestCases)
}

func TestTemplateEscaper_EscapeGitHubActions(t *testing.T) {
	runEscapeTest(
		t,
		"EscapeGitHubActions",
		(*TemplateEscaper).EscapeGitHubActions,
		githubActionsEscapeTestCases,
	)
}

func TestTemplateEscaper_EscapeJSON(t *testing.T) {
	runEscapeTest(t, "EscapeJSON", (*TemplateEscaper).EscapeJSON, jsonEscapeTestCases)
}

func TestTemplateEscaper_EscapeDockerLabel(t *testing.T) {
	runEscapeTest(
		t,
		"EscapeDockerLabel",
		(*TemplateEscaper).EscapeDockerLabel,
		dockerLabelEscapeTestCases,
	)
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
		{
			"YAML injection",
			"name: ${SCRIPT}",
			"yaml",
			false,
		}, // Pattern doesn't include full variable
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
	tests := []valueTestCase[bool]{
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
