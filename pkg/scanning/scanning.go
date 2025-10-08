package scanning

import (
    "fmt"
    "os"

    "github.com/wellatleastitried/yfile/pkg/scanning/matcher"
)

var YaraRulesLoadError = "An unknown error occurred while loading the pre-compiled YARA rules."

// Layer of abstraction for alternate analysis methods to be added in the future
func AnalyzeFile(filePath *string, verbose *bool) {
    rules, err := matcher.LoadEmbeddedRules()
    if err != nil {
        fmt.Fprintln(os.Stderr, YaraRulesLoadError)
    }
    defer rules.Destroy()

    matcher.ShowYaraMatches(filePath, rules, verbose)
}

