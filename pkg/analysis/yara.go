package analysis

import (
    "fmt"

    //"github.com/hillu/go-yara/v4"
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
