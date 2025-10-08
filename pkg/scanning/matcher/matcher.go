package matcher

import (
    "bytes"
    _ "embed"
    "fmt"
    "os"
    "time"

    "github.com/hillu/go-yara/v4"
)

var ErrorScanningFile = "An error occurred while scanning the provided file."

//go:embed ruleset.compiled
var compiledRules []byte

type Callback struct {
    matches []yara.MatchRule
}

func (c *Callback) RuleMatching(_ *yara.ScanContext, rule *yara.Rule) (bool, error) {
    match := yara.MatchRule{
        Rule: rule.Identifier(),
        Namespace: rule.Namespace(),
        Tags: rule.Tags(),
    }
    c.matches = append(c.matches, match)

    return true, nil
}

func (c *Callback) RuleNotMatching(_ *yara.Rule) (bool, error) {
    return true, nil
}

func LoadEmbeddedRules() (*yara.Rules, error) {
    return yara.ReadRules(bytes.NewReader(compiledRules))
}

func ShowYaraMatches(filePath *string, rules *yara.Rules, verbose *bool) {
    defer func() {
        if err := yara.Finalize(); err != nil {
            fmt.Fprintln(os.Stderr, "[Warning] An error occurred while finalizing go-yara: ", err)
        }
    }()

    callback := &Callback{
        matches: make([]yara.MatchRule, 0),
    }

    err := rules.ScanFile(*filePath, 0, 30*time.Second, callback)
    if err != nil {
        fmt.Fprintln(os.Stderr, ErrorScanningFile)
    }

    displayMatches(callback, verbose)
}

func displayMatches(callback *Callback, verbose *bool) {
    if len(callback.matches) == 0 {
        fmt.Println("No matches found.")
        return
    }

    fmt.Printf("%d YARA matches:\n", len(callback.matches))
    for _, match := range callback.matches {
        fmt.Printf("  - Rule: %s (Namespace: %s)\n", match.Rule, match.Namespace)
        if *verbose {
            if len(match.Tags) > 0 {
                fmt.Printf("    Tags: %v\n", match.Tags)
            }

            for _, str := range match.Strings {
                fmt.Printf("    Matched string '%s' at offset %d\n", str.Name, str.Offset)
            }
        }
    }
}

