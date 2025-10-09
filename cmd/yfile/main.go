// Package main is the entry point for the yfile application.
//
// yfile is an extension of the Unix `file` command, designed
// to provide file type detection in addition to malware
// signature detection and analysis.
package main

import (
    "fmt"
    "os"

    "github.com/wellatleastitried/yfile/pkg/argparse"
    "github.com/wellatleastitried/yfile/pkg/linuxfile"
    "github.com/wellatleastitried/yfile/pkg/scanning"
)

// Args should pass through to the `file` command and just have a few addons for `yfile` specific stuff
func main() {
    verbose := argparse.SetBool("v", "verbose", "Enable verbose output", false)
    fileCommandArgs := argparse.SetString("f", "file-args", "Arguments to pass through to the `file` command (e.g. '-b -i')", false)
    fileCommandHelp := argparse.SetBool("fh", "file-help", "Show help for the `file` command and exit", false)

    argparse.Parse()

    if fileCommandHelp {
        displayHelp()
        os.Exit(0)
    }

    files, err := argparse.RetrieveFiles()
    if err != nil {
        fmt.Errorf(err)
        os.Exit(1)
    }

    processFiles(files, fileCommandArgs, verbose)
}

func displayHelp() {
    argparse.printUsage()
}

func runFileCommand(filePath *string, fileCommandArgs *string) {
    if *fileCommandArgs != "" {
        cmd, err := linuxfile.NewCommandWithArgs(filePath, fileCommandArgs)
        if err != nil {
            fmt.Fprintln(os.Stderr, "[Error] (-f, --file-args) flag is invalid:", err)
            os.Exit(1)
        }
        cmd.Execute()
        return
    }

    cmd := linuxfile.NewCommand(filePath)
    cmd.Execute()
}

func processFiles(filePaths []*string, fileCommandArgs *string, verbose *bool) {
    for _, filePath := range filePaths {
        if verifyFilePath(filePath) {
            runFileCommand(filePath, fileCommandArgs)
            scanning.AnalyzeFile(filePath, verbose)
        }
    }
}

func verifyFilePath(filePath *string) bool {
    if _, err := getFileInfo(filePath); err != nil {
        return false
    }

    return true
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

