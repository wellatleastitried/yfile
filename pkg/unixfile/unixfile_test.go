package unixfile

import (
    "testing"

    "github.com/stretchr/testify/require"
)

const ShouldPass = true
const ShouldFail = false

const WantErr = true
const NoErr = false

func TestArgumentParsing(t *testing.T) {
    tests := []struct {
        input    string
        expected bool
    }{
        {"--help", ShouldPass},
        {"--magic-file /tmp/testfile1:/tmp/testfile2", ShouldPass},
        {"--mime-type --keep-going --print0 --list", ShouldPass},
        {"--mime-type ./README.md", ShouldFail},
        {"/tmp/testfile1", ShouldFail},
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            got := commandArgsAreValid(&tt.input)
            require.Equal(t, tt.expected, got, "input: %q", tt.input)
        })
    }
}

