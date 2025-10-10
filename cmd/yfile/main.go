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
    "github.com/wellatleastitried/yfile/pkg/utils"
)

const argParseErrorString = "Error setting up argument parser:"

func main() {
    help, err := argparse.SetBool("h", "help", "Show this help message and exit", false)
    if err != nil {
        fmt.Fprintln(os.Stderr, argParseErrorString, err)
        os.Exit(utils.ExitError)
    }
    verbose, err := argparse.SetBool("v", "verbose", "Enable verbose output", false)
    if err != nil {
        fmt.Fprintln(os.Stderr, argParseErrorString, err)
        os.Exit(utils.ExitError)
    }
    version, err := argparse.SetBool("", "version", "Show version information and exit", false)
    if err != nil {
        fmt.Fprintln(os.Stderr, argParseErrorString, err)
        os.Exit(utils.ExitError)
    }
    //json := argparse.SetBool("j", "json", "Output results in JSON format", false)
    fileCommandArgs, err := argparse.SetString("f", "file-args", "Arguments to pass through to the `file` command (e.g. '-b -i')", false, "")
    if err != nil {
        fmt.Fprintln(os.Stderr, argParseErrorString, err)
        os.Exit(utils.ExitError)
    }
    fileCommandHelp, err := argparse.SetBool("fh", "file-help", "Show help for the `file` command and exit", false)
    if err != nil {
        fmt.Fprintln(os.Stderr, argParseErrorString, err)
        os.Exit(utils.ExitError)
    }

    argparse.Parse()

    if len(os.Args) < 2 || *help {
        displayHelp()
        os.Exit(utils.ExitOk)
    } else if *fileCommandHelp {
        linuxfile.DisplayFileHelp()
        os.Exit(utils.ExitOk)
    } else if *version {
        displayVersion()
        os.Exit(utils.ExitOk)
    }

    files, err := argparse.RetrieveFiles()
    if err != nil {
        fmt.Fprintln(os.Stderr, "No files provided:", err)
        os.Exit(utils.ExitError)
    }

    exitcode := processFiles(files, fileCommandArgs, verbose)
    os.Exit(exitcode)
}

func displayHelp() {
    argparse.PrintUsage()
}

func displayVersion() {
    fmt.Printf("yfile version %s\n", utils.Version)
}

func processFiles(filePaths []string, fileCommandArgs *string, verbose *bool) int {
    for _, filePath := range filePaths {
        linuxfile.RunFileCommand(filePath, fileCommandArgs)

        exitcode := scanning.AnalyzeFile(filePath, verbose)
        if exitcode != utils.ExitOk {
            return exitcode
        }
    }

    return utils.ExitOk
}

