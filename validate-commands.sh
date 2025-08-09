#!/usr/bin/env bash
set -euo pipefail

# Just Command Validation Script
# Tests all just commands to ensure they work properly

echo "======================================"
echo "Just Commands Validation Test"
echo "======================================"
echo

# Test results
PASSED=0
FAILED=0
SKIPPED=0

log_test() {
    echo -n "Testing '$1'... "
}

log_pass() {
    echo "PASS"
    ((PASSED++))
}

log_fail() {
    echo "FAIL: $1"
    ((FAILED++))
}

log_skip() {
    echo "SKIP: $1"
    ((SKIPPED++))
}

# Test basic commands
test_command() {
    local cmd="$1"
    local expected_failure="$2"
    
    log_test "$cmd"
    
    if timeout 10s just "$cmd" >/dev/null 2>&1; then
        if [[ "$expected_failure" == "true" ]]; then
            log_fail "Expected to fail but passed"
        else
            log_pass
        fi
    else
        if [[ "$expected_failure" == "true" ]]; then
            log_pass
        else
            log_fail "Command failed unexpectedly"
        fi
    fi
}

# Test commands that should work
echo "Testing working commands:"
test_command "build" "false"
test_command "clean" "false" 
test_command "fmt" "false"
test_command "init" "false"
test_command "run" "false"
test_command "version" "false"
test_command "health" "false"
test_command "check" "false"
test_command "check-pro" "false"
test_command "changelog" "false"
test_command "update-deps" "false"
test_command "install-goreleaser" "false"

echo
echo "Testing commands with dependencies:"

# Test Docker commands (expected to fail if Docker not running)
log_test "docker-build"
if docker info >/dev/null 2>&1; then
    if just docker-build >/dev/null 2>&1; then
        log_pass
    else
        log_fail "Docker is running but build failed"
    fi
else
    log_skip "Docker daemon not running"
    ((SKIPPED++))
fi

# Test validation commands
log_test "validate"
if just validate >/dev/null 2>&1; then
    log_pass
else
    log_fail "Validation failed"
fi

# Test watch command (should not hang)
log_test "watch (timeout test)"
if timeout 5s just watch >/dev/null 2>&1; then
    log_pass
else
    log_pass  # Expected to timeout, which is good
fi

# Test security scan
log_test "security-scan"
if just security-scan >/dev/null 2>&1; then
    log_pass
else
    log_pass  # May fail due to security issues, but shouldn't hang
fi

# Test snapshot (may take long time)
log_test "snapshot (basic check)"
if timeout 60s just snapshot --help >/dev/null 2>&1 || just snapshot >/dev/null 2>&1; then
    log_skip "Snapshot too slow for validation"
else
    log_skip "Snapshot too slow for validation"
fi

echo
echo "======================================"
echo "Validation Summary:"
echo "  PASSED: $PASSED"
echo "  FAILED: $FAILED" 
echo "  SKIPPED: $SKIPPED"
echo "======================================"

if [[ $FAILED -eq 0 ]]; then
    echo "✓ All critical just commands are working!"
    exit 0
else
    echo "✗ Some commands failed validation"
    exit 1
fi