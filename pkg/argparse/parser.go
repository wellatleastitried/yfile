package argparse

import (
    "errors"
    "fmt"
    "os"

    "github.com/wellatleastitried/yfile/pkg/utils"
)

type Flag struct {
    ShortFormName string
    LongFormName string
    Description string
    Required bool
    reference any
}

type Flags struct {
    getByShortName map[string]Flag
    getByLongName map[string]Flag
}

var flags = &Flags{
    make(map[string]Flag),
    make(map[string]Flag),
}
var defaultValues = make(map[any]string)

var filePathsToProcess = make([]string, 0)
var files = make([]string, 0)

var ErrNoFileProvided = errors.New("no file paths were provided as arguments")
var ErrInvalidFlag = errors.New("invalid flag provided")
var ErrDuplicateFlag = errors.New("a flag with the same short or long form name already exists")

// shortForm can be an empty string if no short form is desired
// longForm can be an empty string if no long form is desired
// HOWEVER, at least one of shortForm or longForm must be provided, otherwise an error is returned
//
// returns a pointer to the bool variable that will hold the flag's value after parsing
func SetBool(shortForm string, longForm string, description string, required bool) (*bool, error) {
    if shortForm == "" && longForm == "" {
        return nil, ErrInvalidFlag
    }

    if shortForm != "" {
        shortForm = "-" + shortForm
    }
    if longForm != "" {
        longForm = "--" + longForm
    }
    flag := Flag{
        ShortFormName: shortForm,
        LongFormName: longForm,
        Description: description,
        Required: required,
        reference: new(bool),
    }

    flags.getByShortName[shortForm] = flag
    flags.getByLongName[longForm] = flag
    return flag.reference.(*bool), nil
}

// shortForm can be an empty string if no short form is desired
// longForm can be an empty string if no long form is desired
// HOWEVER, at least one of shortForm or longForm must be provided, otherwise an error is returned
//
// returns a pointer to the string variable that will hold the flag's value after parsing
func SetString(shortForm string, longForm string, description string, required bool, defaultValue string) (*string, error) {
    if shortForm == "" && longForm == "" {
        return nil, ErrInvalidFlag
    }

    if shortForm != "" {
        shortForm = "-" + shortForm
    }
    if longForm != "" {
        longForm = "--" + longForm
    }

    flag := Flag{
        ShortFormName: shortForm,
        LongFormName: longForm,
        Description: description,
        Required: required,
        reference: new(string),
    }

    err := flags.checkDuplicate(flag)
    if err != nil {
        return nil, err
    }

    defaultValues[flag.reference] = defaultValue
    flags.getByShortName[shortForm] = flag
    flags.getByLongName[longForm] = flag
    return flag.reference.(*string), nil
}

func (f *Flags) contains(arg string) (Flag, bool) {
    if flag, exists := f.getByShortName[arg]; exists {
        return flag, true
    }
    if flag, exists := f.getByLongName[arg]; exists {
        return flag, true
    }

    return Flag{}, false
}

func (f *Flags) checkDuplicate(flag Flag) error {
    if _, result := f.contains(flag.ShortFormName); result {
        return ErrDuplicateFlag
    }
    if _, result := f.contains(flag.LongFormName); result {
        return ErrDuplicateFlag
    }
    return nil
}

func RetrieveFiles(recurse *bool) ([]string, error) {
    for _, file := range filePathsToProcess {
        if utils.IsFile(file) {
            files = append(files, file)
            continue
        } else if utils.IsDir(file) {
            filesFromDir, err := utils.ExtractFilesFromDir(file, recurse)
            if err == nil {
                files = append(files, filesFromDir...)
            }
        }
    }
    if len(files) < 1 {
        return []string{}, ErrNoFileProvided
    }

    return files, nil
}

func Parse() {
    args := os.Args[1:]
    for i := 0; i < len(args); i++ {
        arg := args[i]
        if flag, exists := flags.contains(arg); exists {
            switch v := flag.reference.(type) {
            case *bool:
                *v = true
            case *string:
                if i+1 < len(args) {
                    *v = args[i+1]
                    i++
                } else {
                    fmt.Fprintf(os.Stderr, "[Error] Missing value for flag: %s\n", arg)
                    os.Exit(utils.ExitError)
                }
            }
        } else {
            if utils.IsFile(arg) || utils.IsDir(arg) {
                filePathsToProcess = append(filePathsToProcess, arg)
                continue
            }

            fmt.Fprintf(os.Stderr, "[Error] Unknown argument: %s\n", arg)
            os.Exit(utils.ExitError)
        }
    }
}

func PrintUsage() {
    fmt.Fprintf(os.Stdout, "Usage: %s [options] <file1> <file2> ...\n\nOptions:\n", os.Args[0])
    for _, flag := range flags.getByShortName {
        req := ""
        def := ""
        if flag.Required {
            req = " (required)"
        }
        if defaultValue, exists := defaultValues[flag.ShortFormName]; exists && defaultValue != "" {
            def = fmt.Sprintf(" (default: %s)", defaultValue)
        }
        fmt.Fprintf(os.Stdout, "  %-4s %-20s %s%s%s\n", flag.ShortFormName, flag.LongFormName, flag.Description, req, def)
    }
}

