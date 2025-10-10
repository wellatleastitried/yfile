package utils

import (
    "strings"
)

const (
    FormatText = false
    FormatJSON = true
)

func MaxLineLength(strings ...string) int {
    maxLen := 0
    for _, str := range strings {
        lines := splitLines(str)
        for _, line := range lines {
            if len(line) > maxLen {
                maxLen = len(line)
            }
        }
    }
    return maxLen
}

func splitLines(s string) []string {
    return strings.Split(s, "\n")
}

func ToJSON(fileOutput string, scanOutput string) string {
    // ToJSON not implemented yet
    return fileOutput + "\n" + scanOutput
}

