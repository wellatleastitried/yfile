package analysis

import (
    "fmt"

    "github.com/wellatleastitried/yfile/pkg/analysis/malware"
)

func AnalyzeFile(filePath *string) {
    rules := malware.StageYaraRules()
    fmt.Println("Analyzing file:", *filePath)
}

func (r *Rules) ScanFile(filePath *string) {
    fmt.Println("Scanning file with Yara rules:", *filePath)
    scanner, err := yara.NewScanner(r.CompiledRules)
    if err != nil {
        fmt.Println("Error creating Yara scanner:", err)
        return
    }
    var matches yara.MatchRules
    err = scanner.SetCallback(&matches).ScanFile(*filePath)
    if err != nil {
        fmt.Println("Error scanning file:", err)
        return
    }
    if len(matches) > 0 {
        fmt.Println("Yara matches found:")
        for _, match := range matches {
            fmt.Printf("- Rule: %s\n", match.Rule)
            for _, str := range match.Strings {
                fmt.Printf("  - Offset: %d, Data: %s\n", str.Offset, str.Data)
            }
        }
    } else {
        fmt.Println("No Yara matches found.")
    }
