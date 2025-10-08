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
    verbose := flag.Bool("v", false, "Enable verbose output")
    filePath := flag.String("f", "", "Path to the file to analyze (required)")
    fileCommandArgs := flag.String("file-args", "", "Arguments to pass through to the `file` command (e.g. '-b -i')")

    flag.Parse()

    verifyFilePath(filePath)

    runFileCommand(filePath, fileCommandArgs)

    scanning.AnalyzeFile(filePath, verbose)
}

func runFileCommand(filePath *string, fileCommandArgs *string) {
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
}

// TODO: Make this support multiple files
func verifyFilePath(filePath *string) {
    if *filePath == "" {
        fmt.Fprintln(os.Stderr, "[Error] -file flag is required")
        flag.Usage()
        os.Exit(1)
    }

    if _, err := getFileInfo(filePath); err != nil {
        os.Exit(1)
    }
}

func getFileInfo(filePath *string) (os.FileInfo, error) {
    fileInfo, err := os.Stat(*filePath)
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Fprintf(os.Stderr, "[Error] Path does not exist: %s\n", *filePath)
            return nil, err
        }
        fmt.Fprintf(os.Stderr, "[Error] Could not retrieve file information for %s: %v\n", *filePath, err)
        return nil, err
    }
    return fileInfo, err
}

// TODO: These are for later implementation
// func isFile(filePath *string) (bool, error) {
//     _, err := getFileInfo(filePath)
//     if err != nil {
//         return false, err
//     }
//     return true, nil
// }

// func isDirectory(filePath *string) (bool, error) {
//     fileInfo, err := getFileInfo(filePath)
//     if err != nil {
//         return false, err
//     }
//     return fileInfo.IsDir(), nil
// }

