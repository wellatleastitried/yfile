package linuxfile

import (
    "testing"
)

func TestArgumentParsing(t *testing.T) {
    tests := []struct {
        input    string
        expected []string
        wantErr  bool
    }{
        {"--help", []string{"--help"}, false},
        {"--magic-file /tmp/testfile1:/tmp/testfile2", []string{"--magic-file", "/tmp/testfile1:/tmp/testfile2"}, false},
        {"--mime-type --keep-going --print0 --list", []string{"--mime-type", "--keep-going", "--print0", "--list"}, false},
        {"", []string{}, false},
    }

    for _, tt := range tests {
        arg, err := parseCommandArgs(&tt.input)
        if (err != nil) != tt.wantErr {
            t.Errorf("ParseArgument(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
            continue
        }
        if err == nil && stringSlicesEqual(arg, tt.expected) == false {
            t.Errorf("ParseArgument(%q) = %v, want %v", tt.input, arg, tt.expected)
        }
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
