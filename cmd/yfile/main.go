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
    "github.com/wellatleastitried/yfile/pkg/scanning"
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

    scanning.AnalyzeFile(filePath)
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

func isFile(filePath *string) (bool, error) {
    _, err := getFileInfo(filePath)
    if err != nil {
        return false, err
    }
    return true, nil
}

func isDirectory(filePath *string) (bool, error) {
    fileInfo, err := getFileInfo(filePath)
    if err != nil {
        return false, err
    }
    return fileInfo.IsDir(), nil
}

func getFileInfo(filePath *string) (os.FileInfo, error) {
    fileInfo, err := os.Stat(*filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("Path does not exist: %w", *filePath)
        }
        return nil, fmt.Errorf("Error retrieving file information for %w", *filePath)
    }
    return fileInfo, err
}

