package scanning

import (
    "embed"
    "github.com/hillu/go-yara/v4"
)

//go:embed ruleset.compiled
var compiledRules []byte

func loadEmbeddedRules() (*yara.Rules, error) {
    return yara.LoadRulesFromBytes(compiledRules)
}

func showYaraMatches(filePath *string) {
    matches, err := rules.ScanFile(*filePath, 0, 0)
    for _, match := range matches {
        fmt.Println("Matched rule: ", match.Rule)
    }
}

