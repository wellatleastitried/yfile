package unixfile

import (
    "fmt"
    "errors"
    "strings"
    "os"
    "os/exec"

    "github.com/wellatleastitried/yfile/pkg/utils"
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

func DisplayFileHelp() {
    cmd := exec.Command("file", "--help")
    stdout, err := cmd.Output()
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v", err)
        os.Exit(utils.ExitError)
    }

    fmt.Fprintf(os.Stdout, "%s", string(stdout))
}

func newCommand(filePath string) CommandData {
    return CommandData {
        Name: "file",
        FilePath: filePath,
    }
}

func newCommandWithArgs(filePath string, commandArgs *string) (CommandDataWithArgs, error) {
    if !commandArgsAreValid(commandArgs) {
        return CommandDataWithArgs{}, ErrFilePathProvidedInArguments
    }

    return CommandDataWithArgs {
        Name: "file",
        Args: *commandArgs,
        FilePath: filePath,
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
    

func (c CommandData) Execute() string {
    cmd := exec.Command(c.Name, c.FilePath)
    return execute(cmd)
}

func (c CommandDataWithArgs) Execute() string {
    if c.Args == "" {
        fmt.Fprintln(os.Stderr, ErrNoArgumentsProvidedToFileArgsFlag)
        os.Exit(utils.ExitError)
    }

    cmd := exec.Command(c.Name, c.Args, c.FilePath)
    return execute(cmd)
}

func execute(cmd *exec.Cmd) string {
    stdout, err := cmd.Output()
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v", err)
        os.Exit(utils.ExitError)
    }

    return string(stdout)
}

func RunFileCommand(filePath string, fileCommandArgs *string) string {
    if *fileCommandArgs != "" {
        cmd, err := newCommandWithArgs(filePath, fileCommandArgs)
        if err != nil {
            fmt.Fprintln(os.Stderr, "[Error] (-f, --file-args) flag is invalid:", err)
            os.Exit(utils.ExitError)
        }
        return cmd.Execute()
    }

    cmd := newCommand(filePath)
    return cmd.Execute()
}

