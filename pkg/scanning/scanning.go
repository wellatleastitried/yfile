package scanning

import (
    "errors"
    "fmt"
    "os"

    "github.com/wellatleastitried/yfile/pkg/scanning/matcher"
    "github.com/wellatleastitried/yfile/pkg/utils"
)

var ErrYaraRulesLoad = errors.New("unknown error occurred while loading the pre-compiled YARA rules")

// Layer of abstraction for alternate analysis methods to be added in the future
func AnalyzeFile(filePath string, verbose *bool) (string, int) {
    rules, err := matcher.LoadEmbeddedRules()
    if err != nil {
        fmt.Fprintln(os.Stderr, "[Error] ", ErrYaraRulesLoad, ":", err)
    }
    defer rules.Destroy()

    output, count := matcher.ShowYaraMatches(filePath, rules, verbose)
    return output, utils.GetExitcodeFromMatches(count)
}

