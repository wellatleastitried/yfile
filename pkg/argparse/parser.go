package argparse

import (
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

var flags *Flags = &Flags{make(map[string]Flag)}
var defaultValues map[string]string = make(map[string]string)

var files []string = make([]string, 0)

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
				for _, file := range filesFromDir {
					files = append(files, file)
				}
				continue
			}
			fmt.Errorf("Unknown argument: %s", arg)
			os.Exit(1)
        }
    }
}

