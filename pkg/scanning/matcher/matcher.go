package matcher

import (
    "bytes"
    _ "embed"
    "errors"
    "fmt"
    "os"
    "time"

    "github.com/hillu/go-yara/v4"
)

var ErrorScanningFile = errors.New("error occurred while scanning the provided file")

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

    return false, nil
}

func (c *Callback) RuleNotMatching(_ *yara.Rule) (bool, error) {
    return false, nil
}

func LoadEmbeddedRules() (*yara.Rules, error) {
    return yara.ReadRules(bytes.NewReader(compiledRules))
}

func ShowYaraMatches(filePath string, rules *yara.Rules, verbose *bool) (string, int) {
    callback := &Callback{
        matches: make([]yara.MatchRule, 0),
    }

    err := rules.ScanFile(filePath, 0, 30*time.Second, callback)
    if err != nil {
        fmt.Fprintln(os.Stderr, ErrorScanningFile)
    }

    output := getMatches(callback, verbose)

    return output, len(callback.matches)
}

func getMatches(callback *Callback, verbose *bool) string {
    if len(callback.matches) == 0 {
        return "File does not match common malware signatures."
    }

    output := fmt.Sprintf("%d YARA matches:\n", len(callback.matches))
    for _, match := range callback.matches {
        output += fmt.Sprintf("  - Rule: %s (Namespace: %s)\n", match.Rule, match.Namespace)
        if *verbose {
            if len(match.Metas) > 0 {
                output += fmt.Sprintf("    Metas: %v\n", match.Metas)
            }
            if len(match.Tags) > 0 {
                output += fmt.Sprintf("    Tags: %v\n", match.Tags)
            }

            for _, str := range match.Strings {
                fmt.Printf("    Matched string '%s' at offset %d\n", str.Name, str.Offset)
                output += fmt.Sprintf("    Matched string '%s' at offset %d\n", str.Name, str.Offset)
            }
        }
    }

    return output
}

