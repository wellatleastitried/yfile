#!/usr/bin/env bash

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

BINARY="./build/yfile"
FAILED=0
PASSED=0

test_exit_code() {
    local test_name="$1"
    local expected_exit="$2"
    shift 2
    local cmd=("$@")

    echo -n "Testing: $test_name - "

    set +e
    "${cmd[@]}" > /dev/null 2>&1
    actual_exit=$?
    set -e

    if [ "$actual_exit" -eq "$expected_exit" ]; then
        echo -e "${GREEN}PASS${NC}"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}FAILED${NC}"
        FAILED=$((FAILED + 1))
    fi
}

test_output() {
    local test_name="$1"
    local expected_pattern="$2"
    shift 2
    local cmd=("$@")

    echo -n "Testing: $test_name - "

    output=$("${cmd[@]}" 2>&1 || true)

    if echo "$output" | grep -qi "$expected_pattern"; then
        echo -e "${GREEN}PASS${NC}"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}FAILED${NC}"
        FAILED=$((FAILED + 1))
    fi
}

TEST_DIR=$(mktemp -d)
trap "rm -rf $TEST_DIR" EXIT

echo "Test data" > "$TEST_DIR/test.txt"
echo "#!/bin/bash" > "$TEST_DIR/test.sh"
chmod +x "$TEST_DIR/test.sh"

echo "Running integration tests..."
echo

test_exit_code "Valid text file" 0 "$BINARY" "$TEST_DIR/test.txt"
test_exit_code "Valid script file" 0 "$BINARY" "$TEST_DIR/test.sh"
test_exit_code "Non-existent file" 2 "$BINARY" "/nonexistent/file.txt"
test_exit_code "No arguments" 0 "$BINARY"
test_exit_code "Multiple files" 0 "$BINARY" "$TEST_DIR/test.txt" "$TEST_DIR/test.sh"
test_exit_code "Verbose mode" 0 "$BINARY" "-v" "$TEST_DIR/test.txt"

test_output "Output contains filename" "test.txt" "$BINARY" "$TEST_DIR/test.txt"
test_output "Text file detection" "text" "$BINARY" "$TEST_DIR/test.txt"
test_output "Script file detection" "script" "$BINARY" "$TEST_DIR/test.sh"

echo
echo "================================"
echo -e "Tests passed: ${GREEN}$PASSED${NC}"
echo -e "Tests failed: ${RED}$FAILED${NC}"
echo "================================"

if [ "$FAILED" -gt 0 ]; then
    exit 1
fi

