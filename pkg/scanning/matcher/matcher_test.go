package matcher

import (
    "bytes"
    "fmt"
    "io"
    "os"
    "strings"
    "testing"

    "github.com/stretchr/testify/require"
)

const ShouldPass = true
const ShouldFail = false

const WantErr = true
const NoErr = false

func TestLoadEmbeddedRules(t *testing.T) {
    rules, err := LoadEmbeddedRules()
    if err != nil {
        t.Fatalf("LoadEmbeddedRules() failed: %v", err)
    }

    if rules == nil {
        t.Fatal("LoadEmbeddedRules() returned nil rules")
    }

    // Verify rules can be used
    defer rules.Destroy()
}

func TestLoadEmbeddedRulesMultipleTimes(t *testing.T) {
    for i := 0; i < 5; i++ {
        rules, err := LoadEmbeddedRules()
        if err != nil {
            t.Fatalf("LoadEmbeddedRules() failed on iteration %d: %v", i, err)
        }

        if rules == nil {
            t.Fatalf("LoadEmbeddedRules() returned nil rules on iteration %d", i)
        }

        rules.Destroy()
    }
}

func TestShowYaraMatches(t *testing.T) {
    filecontents := [][]byte{
        []byte("This is a safe file"),
        []byte("LUA_PATH Hi. Happy reversing, you can mail me: luabot@yandex.ru /tmp/lua_XXXXXX NOTIFY UPDATE"),
        []byte("Another safe file."),
        []byte("<?php $tkqagunjoiapdiytpnmxuthqyhqbadqargdvlv = ' $kv=explode(chr((21726423342181308796293439586547227+9948363247998)) $cegubx=(05966373430082365026859-53816062886690162) if (!function_exists('znkbrbixfqtsfhhcbjywjasxmsjyvafi'))"),
    }

    for i, content := range filecontents {
        fileName := fmt.Sprintf("testfile_%d.txt", i)
        tmpFile, err := os.CreateTemp("", fileName)
        if err != nil {
            t.Fatalf("Iteration %d:\nFailed to create temp file: %v", i, err)
        }
        defer os.Remove(tmpFile.Name())

        if _, err := tmpFile.Write(content); err != nil {
            t.Fatalf("Iteration %d:\nFailed to write to temp file: %v", i, err)
        }
        tmpFile.Close()

        rules, err := LoadEmbeddedRules()
        if err != nil {
            t.Fatalf("Iteration %d:\nFailed to load embedded rules: %v", i, err)
        }
        defer rules.Destroy()

        filePath := tmpFile.Name()
        verbose := false

        oldStdout := os.Stdout
        oldStderr := os.Stderr
        r, w, _ := os.Pipe()
        os.Stdout = w
        os.Stderr = w

        ShowYaraMatches(&filePath, rules, &verbose)

        w.Close()
        os.Stdout = oldStdout
        os.Stderr = oldStderr

        var buf bytes.Buffer
        if _, err := io.Copy(&buf, r); err != nil {
            t.Fatalf("Iteration %d:\nFailed to read captured output: %v", i, err)
        }

        output := buf.String()

        if output == "" {
            t.Fatalf("Iteration %d\nNo output captured from ShowYaraMatches()", i)
        } else if strings.Contains(output, "File does not match common malware signatures.") {
            if bytes.Contains(content, []byte("safe")) {
                // Expected no matches
                continue
            }
            t.Fatalf("Iteration %d\nExpected matches but got none. Output: %s", i, output)
        }

        if bytes.Contains(content, []byte("?php")) || bytes.Contains(content, []byte("LUA_PATH")) {
            require.Contains(t, output, "Rule:")
        }
    }
}

