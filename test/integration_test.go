//go:build integration
// +build integration

package test

import (
    "fmt"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
)

func TestFileIntegration(t *testing.T) {
    cmd := exec.Command("pwd")
    output, _ := cmd.CombinedOutput()
    fmt.Println(string(output))

    binaryPath := filepath.Join("..", "build", "yfile")

    tests := []struct {
        name string
        file string
        args []string
        wantExitCode int
        wantInOutput string
    }{
        {
            name: "clean text file",
            file: "testsignatures/clean.txt",
            args: []string{},
            wantExitCode: 0,
            wantInOutput: "File does not match common malware signatures",
        },
        {
            name: "malicious lua file",
            file: "testsignatures/lua_malware_sig.lua",
            args: []string{},
            wantExitCode: 1,
            wantInOutput: "Rule: LuaBot",
        },
        {
            name: "misc malicious file",
            file: "testsignatures/multiple.exe.txt",
            args: []string{},
            wantExitCode: 1,
            wantInOutput: "3 YARA matches:",
        },
        {
            name: "non-existent file",
            file: "testsignatures/nonexistent.file",
            args: []string{},
            wantExitCode: 2,
            wantInOutput: "",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            args := append(tt.args, tt.file)
            cmd := exec.Command(binaryPath, args...)
            output, err := cmd.CombinedOutput()

            exitCode := 0
            if err != nil {
                if exitErr, ok := err.(*exec.ExitError); ok {
                    exitCode = exitErr.ExitCode()
                } else {
                    t.Fatalf("Command failed: %v", err)
                }
            }

            if exitCode != tt.wantExitCode {
                t.Errorf("got exit code %d, want %d", exitCode, tt.wantExitCode)
            }

            if !strings.Contains(string(output), tt.wantInOutput) {
                t.Errorf("output does not contain expected string.\ngot: %s\nwant to contain: %s", string(output), tt.wantInOutput)
            }
        })
    }
}
