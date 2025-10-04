// Package main is the entry point for the yfile application.
//
// yfile is an extension of the Unix `file` command, designed
// to provide file type detection in addition to malware
// signature detection and analysis.
package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/wellatleastitried/yfile/pkg/linuxfile"
)

// Args should pass through to the `file` command and just have a few addons for `yfile` specific stuff
func main() {
    filePath := flag.String("f", "", "Path to the file to analyze (required)")
    fileCommandArgs := flag.String("file-args", "", "Arguments to pass through to the `file` command (e.g. '-b -i')")

    flag.Parse()

    verifyFilePath(filePath)

    if *fileCommandArgs != "" {
        cmd, err := linuxfile.NewCommandWithArgs(filePath, fileCommandArgs)
        if err != nil {
            fmt.Fprintln(os.Stderr, "[Error] -file-args flag is invalid:", err)
            os.Exit(1)
        }
        cmd.Execute()
    } else {
        cmd := linuxfile.NewCommand(filePath)
        cmd.Execute()
    }

    // TODO: Replace this with actually kicking off detection and analysis
}

// TODO: Make this support multiple files
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

