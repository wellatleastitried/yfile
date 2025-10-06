package analysis

import (
    "fmt"

    "github.com/hillu/go-yara/v4"
)

const (
    // Change to something more permanent like /usr/local/share/yara or similar
    YaraRulesPath = "./rules/index_all.yar"
)

type Match struct {
    Rule string
    Namespace string
    Tags []string
    Strings []string
}

type Result struct {
    Matches []Match
    Error error
}

// Use Yara rules to detect malware signatures
func ImportYaraRules() {
    fmt.Println("Importing Yara rules...")
}

