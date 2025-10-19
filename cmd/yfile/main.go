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
    "sync"

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
    recurse := setBoolFlag("r", "recurse", "Recurse into directories when provided as input", false)
    threads := setStringFlag("t", "threads", "Number of threads to use for processing files (min 1, default 6, max " + utils.CPUCountString() + ")", false, "8")
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

    exitcode := processFiles(files, fileCommandArgs, threads, verbose)
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

// fileResult holds the output and exit code for a single file analysis
type fileResult struct {
    filePath    string
    fileOutput  string
    scanOutput  string
    exitcode    int
}

func processFiles(filePaths []string, fileCommandArgs, numWorkers *string, verbose *bool) int {
    exitcode := utils.ExitOk

    workers := utils.GetWorkerCount(numWorkers)

    jobs := make(chan string, len(filePaths))
    results := make(chan fileResult, len(filePaths))
    var wg sync.WaitGroup
    
    for w := 0; w < workers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for path := range jobs {
                fileOutput := unixfile.RunFileCommand(path, fileCommandArgs)
                scanOutput, currentExitcode := scanning.AnalyzeFile(path, verbose)
                
                results <- fileResult{
                    filePath:   path,
                    fileOutput: fileOutput,
                    scanOutput: scanOutput,
                    exitcode:   currentExitcode,
                }
            }
        }()
    }
    
    for _, filePath := range filePaths {
        jobs <- filePath
    }
    close(jobs)
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    for result := range results {
        if result.exitcode == utils.ExitError {
            return result.exitcode
        }
        
        if result.exitcode == utils.ExitInfected {
            // Set exit code to infected if ANY single file is infected
            exitcode = utils.ExitInfected
        }
        
        displayOutput(result.fileOutput, result.scanOutput)
    }
    
    return exitcode
}

func displayOutput(fileOutput, scanOutput string) {
    fmt.Printf("%s%s\n", fileOutput, scanOutput)

    maxLen := utils.MaxLineLength(fileOutput, scanOutput)
    if maxLen > 0 {
        fmt.Println(strings.Repeat("-", maxLen))
    }
}

