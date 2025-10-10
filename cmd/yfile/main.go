// Package main is the entry point for the yfile application.
//
// yfile is an extension of the Unix `file` command, designed
// to provide file type detection in addition to malware
// signature detection and analysis.
package main

import (
    "fmt"
    "os"
    "strings"

    "github.com/wellatleastitried/yfile/pkg/argparse"
    "github.com/wellatleastitried/yfile/pkg/unixfile"
    "github.com/wellatleastitried/yfile/pkg/scanning"
    "github.com/wellatleastitried/yfile/pkg/utils"
)

const argParseErrorString = "Error setting up argument parser:"

func main() {
    help := setBoolFlag("h", "help", "Show this help message and exit", false)
    verbose := setBoolFlag("v", "verbose", "Enable verbose output", false)
    version := setBoolFlag("", "version", "Show version information and exit", false)
    json := setBoolFlag("j", "json", "Output results in JSON format", false)
    recurse := setBoolFlag("r", "recurse", "Recurse into directories when provided as input", false)
    fileCommandArgs := setStringFlag("f", "file-args", "Arguments to pass through to the `file` command (e.g. '-b -i')", false, "")
    fileCommandHelp := setBoolFlag("fh", "file-help", "Show help for the `file` command and exit", false)

    argparse.Parse()

    if len(os.Args) < 2 || *help {
        displayHelp()
        os.Exit(utils.ExitOk)
    } else if *fileCommandHelp {
        unixfile.DisplayFileHelp()
        os.Exit(utils.ExitOk)
    } else if *version {
        displayVersion()
        os.Exit(utils.ExitOk)
    }

    files, err := argparse.RetrieveFiles(recurse)
    if err != nil {
        fmt.Fprintln(os.Stderr, "No files provided:", err)
        os.Exit(utils.ExitError)
    }

    exitcode := processFiles(files, fileCommandArgs, verbose, json)
    os.Exit(exitcode)
}

// I am handling errors the same way for every flag so I wrapped argparse calls
func setStringFlag(shortForm, longForm, description string, required bool, defaultValue string) *string {
    flag, err := argparse.SetString(shortForm, longForm, description, required, defaultValue)
    if err != nil {
        fmt.Fprintln(os.Stderr, argParseErrorString, err)
        os.Exit(utils.ExitError)
    }
    return flag
}

func setBoolFlag(shortForm, longForm, description string, required bool) *bool {
    flag, err := argparse.SetBool(shortForm, longForm, description, required)
    if err != nil {
        fmt.Fprintln(os.Stderr, argParseErrorString, err)
        os.Exit(utils.ExitError)
    }
    return flag
}

func displayHelp() {
    argparse.PrintUsage()
}

func displayVersion() {
    fmt.Printf("yfile version %s\n", utils.Version)
}

func processFiles(filePaths []string, fileCommandArgs *string, verbose , outputFormat *bool) int {
    exitcode := utils.ExitOk
    for _, filePath := range filePaths {
        fileOutput := unixfile.RunFileCommand(filePath, fileCommandArgs)

        scanOutput, exitcode := scanning.AnalyzeFile(filePath, verbose)
        if exitcode == utils.ExitInfected {
            // Set exit code to infected if ANY single file is infected
            exitcode = utils.ExitInfected
        }
        if exitcode == utils.ExitError {
            return exitcode
        }

        displayOutput(outputFormat, fileOutput, scanOutput)
        displayDivider(fileOutput, scanOutput)
    }

    return exitcode
}

func displayOutput(outputFormat *bool, fileOutput, scanOutput string) {
    if *outputFormat == utils.FormatJSON {
        // json output not yet implemented
        fmt.Println(utils.ToJSON(fileOutput, scanOutput))
    } else {
            fmt.Printf("%s%s\n", fileOutput, scanOutput)
    }
}

func displayDivider(fileOutput, scanOutput string) {
    maxLen := utils.MaxLineLength(fileOutput, scanOutput)
    if maxLen > 0 {
        fmt.Println(strings.Repeat("-", maxLen))
    }
}

