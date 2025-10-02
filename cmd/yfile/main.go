package main

import (
	"flag"
	"os"
	"fmt"
)

// Parse args and kick off processes
func main() {
    filePath := flag.String("file", "", "Path to the file to analyze (required)")

    flag.Parse()

    verifyFilePath(filePath)

    // TODO: Replace this with actually kicking off detection and analysis
    fmt.Fprintf(os.Stdout, "[Success] Your file is '%s'\n", *filePath)
}

func verifyFilePath(filePath *string) {
    if *filePath == "" {
        fmt.Fprintln(os.Stderr, "[Error] -file flag is required")
        flag.Usage()
        os.Exit(1)
    }

    if _, err := os.Stat(*filePath); os.IsNotExist(err) {
        fmt.Fprintf(os.Stderr, "[Error] file '%s' does not exist\n", *filePath)
        os.Exit(1)
    }
}

