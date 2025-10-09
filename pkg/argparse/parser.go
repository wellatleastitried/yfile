package argparse

import (
    "errors"
    "fmt"
    "os"
)

type Flag struct {
    Name string
    AlternateName string
    Description string
    Required bool
    reference any
}

type Flags struct {
    stored map[string]Flag
}

var flags = &Flags{make(map[string]Flag)}
var defaultValues = make(map[string]string)

var files = make([]string, 0)

var ErrNoFileProvided = errors.New("no file paths were provided as arguments")

func SetBool(name string, alternateName string, description string, required bool) *bool {
    flag := Flag{
        Name: name,
        AlternateName: alternateName,
        Description: description,
        Required: required,
        reference: new(bool),
    }
    flags.stored[name] = flag
    return flag.reference.(*bool)
}

func SetString(name string, alternateName string, description string, required bool, defaultValue string) *string {
    flag := Flag{
        Name: name,
        AlternateName: alternateName,
        Description: description,
        Required: required,
        reference: new(string),
    }
    defaultValues[name] = defaultValue
    flags.stored[name] = flag
    return flag.reference.(*string)
}

func contains(f *Flags, arg string) (Flag, bool) {
    shortForm := "-" + arg
    longForm := "--" + arg

    if flag, exists := f.stored[shortForm]; exists {
        return flag, true
    }

    if flag, exists := f.stored[longForm]; exists {
        return flag, true
    }

    if flag, exists := f.stored[arg]; exists {
        return flag, true
    }

    return Flag{}, false
}

func RetrieveFiles() ([]string, error) {
    if len(files) < 1 {
        return []string{}, ErrNoFileProvided
    }

    return files, nil
}

func Parse() {
    flags.Parse()
}

func (f *Flags) Parse() {
    args := os.Args[1:]
    for i := 0; i < len(args); i++ {
        arg := args[i]
        if flag, exists := contains(f, arg); exists {
            switch v := flag.reference.(type) {
            case *bool:
                *v = true
            case *string:
                if i+1 < len(args) {
                    *v = args[i+1]
                    i++
                }
            }
        } else {
            if isFile(arg) {
                files = append(files, arg)
                continue
            } else if isDir(arg) {
                filesFromDir := extractFilesFromDir(arg)
                files = append(files, filesFromDir...)
                continue
            }
            fmt.Fprintf(os.Stderr, "[Error] Unknown argument: %s\n", arg)
            os.Exit(1)
        }
    }
}

func PrintUsage() {
    flags.PrintUsage()
}

func (f *Flags) PrintUsage() {
    // TODO: Poll flags and display in a beautiful way
}

