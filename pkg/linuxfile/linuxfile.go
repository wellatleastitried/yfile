package linuxfile

import (
    "encoding/csv"
    "fmt"
    "strings"
    "errors"
    "os"
)

type Command struct {
    Name string
    Args []string
    FilePath string
}

var ArgumentParsingError = errors.New("Unable to parse arguments provided to `file` command")

func NewCommand(filePath *string, commandArgs *string) Command {
    args, err := parseCommandArgs(commandArgs)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    command := Command {
        Name: "file",
        Args: args,
        FilePath: *filePath,
    }

    return command
}

func parseCommandArgs(commandArgs *string) ([]string, error) {
    if *commandArgs == "" {
        return []string{}, nil
    }

    r := csv.NewReader(strings.NewReader(*commandArgs))
    r.Comma = ' '
    fields, err := r.Read()
    if err != nil {
        return []string{}, ArgumentParsingError
    }

    return fields, nil
}

func (c Command) Execute() {
    fmt.Println("Executing command...")
}

