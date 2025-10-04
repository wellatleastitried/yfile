package linuxfile

import (
    "testing"

    "github.com/stretchr/testify/require"
)

const SHOULD_PASS = true
const SHOULD_FAIL = false

const WANT_ERR = true
const NO_ERR = false

func TestArgumentParsing(t *testing.T) {
    tests := []struct {
        input    string
        expected bool
    }{
        {"--help", SHOULD_PASS},
        {"--magic-file /tmp/testfile1:/tmp/testfile2", SHOULD_PASS},
        {"--mime-type --keep-going --print0 --list", SHOULD_PASS},
        {"--mime-type ./README.md", SHOULD_FAIL},
        {"/tmp/testfile1", SHOULD_FAIL},
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            got := commandArgsAreValid(&tt.input)
            require.Equal(t, tt.expected, got, "input: %q", tt.input)
        })
    }
}

func stringSlicesEqual(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}
