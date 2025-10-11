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

    echo -n "Testing: $test_name"

    set +e
    "${cmd[@]}" > /dev/null 2>&1
    actual_exit=$?
    set -e

    if [ "$actual_exit" -eq "$expected_exit" ]; then
        echo -e "${GREEN}PASS${NC}"
        ((PASSED++))
    else
        echo -e "${RED}FAILED${NC}"
        ((FAILED++))
    fi
}

test_output() {
    local test_name="$1"
    local expected_pattern="$2"
    shift 2
    local cmd=("$@")

    echo -n "Testing: $test_name"

    output=$("${cmd[@]}" 2>&1 || true)

    if echo "$output" | grep -qi "$expected_pattern"; then
        echo -e "${GREEN}PASS${NC}"
        ((PASSED++))
    else
        echo -e "${RED}FAILED${NC}"
        ((FAILED++))
    fi
}

TEST_DIR=$(mktemp -d)
trap "rm -rf $TEST_DIR" EXIT

echo "Test data" > "$TEST_DIR/test.txt"
echo "#!/bin/bash" > "$TEST_DIR/test.sh"
chmod +x "$TEST_DIR/test.sh"

echo "Running integration tests..."
echo

# TODO
test_exit_code
test_exit_code
test_exit_code
test_exit_code
test_exit_code
test_exit_code

test_output
test_output
test_output

echo
echo "================================"
echo -e "Tests passed: ${GREEN}$PASSED${NC}"
echo -e "Tests failed: ${RED}$FAILED${NC}"
echo "================================"

if [ "$FAILED" -gt 0 ]; then
    exit 1
fi

