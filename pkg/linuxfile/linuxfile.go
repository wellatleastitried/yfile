package linuxfile

import (
    "fmt"
    "errors"
    "strings"
    "os"
    "os/exec"
)

type CommandData struct {
    Name string
    FilePath string
}

type CommandDataWithArgs struct {
    Name string
    Args string
    FilePath string
}

var ErrNoArgumentsProvidedToFileArgsFlag = errors.New("the `file-args` flag was provided but no arguments were passed to it")
var ErrFilePathProvidedInArguments = errors.New("in yfile the file path can only be provided with the -f flag it cannot reside within the -file-args")

func NewCommand(filePath *string) CommandData {
    return CommandData {
        Name: "file",
        FilePath: *filePath,
    }
}

func NewCommandWithArgs(filePath *string, commandArgs *string) (CommandDataWithArgs, error) {
    if !commandArgsAreValid(commandArgs) {
        return CommandDataWithArgs{}, ErrFilePathProvidedInArguments
    }

    return CommandDataWithArgs {
        Name: "file",
        Args: *commandArgs,
        FilePath: *filePath,
    }, nil
}

func commandArgsAreValid(commandArgs *string) bool {
    args := strings.Split(*commandArgs, " ")
    fileExpected := false
    for _, arg := range args {
        if fileExpected {
            fileExpected = false
            continue
        }
        if strings.HasPrefix(arg, "-m") || strings.Contains(arg, "--magic-file") || strings.Contains(arg, "-f") || strings.Contains(arg, "--files-from") {
            fileExpected = true
            continue
        }
        if !strings.HasPrefix(arg, "-") || strings.Contains(arg, "/") || strings.Contains(arg, ".") {
            return false
        }
    }
    return true
}
    

func (c CommandData) Execute() {
    cmd := exec.Command(c.Name, c.FilePath)
    execute(cmd)
}

func (c CommandDataWithArgs) Execute() {
    if c.Args == "" {
        fmt.Fprintln(os.Stderr, ErrNoArgumentsProvidedToFileArgsFlag)
        os.Exit(1)
    }

    cmd := exec.Command(c.Name, c.Args, c.FilePath)
    execute(cmd)
}

func execute(cmd *exec.Cmd) {
    stdout, err := cmd.Output()
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v", err)
        os.Exit(1)
    }

    fmt.Fprintf(os.Stdout, "%s", string(stdout))
}

