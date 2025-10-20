package main

import (
    "fmt"
    "os"
    "strings"
    "sync"

    "github.com/wellatleastitried/yfile/pkg/argparse"
    "github.com/wellatleastitried/yfile/pkg/scanning"
    "github.com/wellatleastitried/yfile/pkg/unixfile"
    "github.com/wellatleastitried/yfile/pkg/utils"
)

const argParseErrorString = "Error setting up argument parser:"

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

type fileResult struct {
    filePath    string
    fileOutput  string
    scanOutput  string
    exitcode    int
}

func processFiles(filePaths []string, fileCommandArgs, numWorkers *string, verbose *bool) int {
    resultChannel := make(chan fileResult, len(filePaths))
    dispatchScanningThreads(&resultChannel, numWorkers, filePaths, fileCommandArgs, verbose)
    
    exitcode := handleResults(resultChannel)
    return exitcode
}

func dispatchScanningThreads(results *chan fileResult, numWorkers *string, filePaths []string, fileCommandArgs *string, verbose *bool) {
    jobs := make(chan string, len(filePaths))

    workers := utils.GetWorkerCount(numWorkers)
    var wg sync.WaitGroup
    
    for w := 0; w < workers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for path := range jobs {
                fileOutput := unixfile.RunFileCommand(path, fileCommandArgs)
                scanOutput, currentExitcode := scanning.AnalyzeFile(path, verbose)
                
                *results <- fileResult{
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
        close(*results)
    }()
}

func handleResults(results chan fileResult) int {
    exitcode := utils.ExitOk

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

