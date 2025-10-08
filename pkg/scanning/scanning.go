package scanning

// Layer of abstraction for alternate analysis methods to be added in the future
func AnalyzeFile(filePath *string) {
    yara.showYaraMatches(filePath)
}

