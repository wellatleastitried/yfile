// Package main is the entry point for the yfile application.
//
// yfile is an extension of the Unix `file` command, designed
// to provide file type detection in addition to malware
// signature detection and analysis.
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/wellatleastitried/yfile/pkg/argparse"
    "github.com/wellatleastitried/yfile/pkg/unixfile"
    "github.com/wellatleastitried/yfile/pkg/utils"
)

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

    startTime := time.Now()

    files, err := argparse.RetrieveFiles(recurse)
    if err != nil {
        fmt.Fprintln(os.Stderr, "No files provided:", err)
        os.Exit(utils.ExitError)
    }

    exitcode := processFiles(files, fileCommandArgs, threads, verbose)

    endTime := time.Now()
    elapsedTime := endTime.Sub(startTime)
    if *verbose {
        fmt.Printf("File(s) scanned in %s\n", elapsedTime)
    }

    os.Exit(exitcode)
}

