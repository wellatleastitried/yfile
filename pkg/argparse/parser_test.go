package argparse

import (
    "os"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestSetBool(t *testing.T) {
    // Clear previous flags
    flags.getByShortName = make(map[string]Flag)
    flags.getByLongName = make(map[string]Flag)

    errFlag, err := SetBool("", "", "An invalid flag", false)
    require.Nil(t, errFlag, "SetBool should return nil for invalid flag")
    require.Equal(t, ErrInvalidFlag, err, "SetBool should return ErrInvalidFlag for invalid flag")

    boolFlag, _ := SetBool("b", "bool", "A boolean flag", true)
    
    require.NotNil(t, boolFlag, "SetBool should return a non-nil pointer")
    require.Equal(t, false, *boolFlag, "Default value of a bool flag should be false")

    flag, exists := flags.getByShortName["-b"]
    require.True(t, exists, "Flag should be getByShortName in the flags map")
    require.Equal(t, "-b", flag.ShortFormName, "Short form name should match")
    require.Equal(t, "--bool", flag.LongFormName, "Long form name should match")
    require.Equal(t, "A boolean flag", flag.Description, "Description should match")
    require.Equal(t, true, flag.Required, "Required status should match")

    setMyOwnArgs("-b")
    Parse()

    require.Equal(t, true, *boolFlag, "Bool flag should be set to true after parsing")
}

func TestSetString(t *testing.T) {
    // Clear previous flags
    flags.getByLongName = make(map[string]Flag)
    defaultValues = make(map[any]string)

    errFlag, err := SetString("", "", "An invalid flag", false, "default")
    require.Nil(t, errFlag, "SetString should return nil for invalid flag")
    require.Equal(t, ErrInvalidFlag, err, "SetString should return ErrInvalidFlag for invalid flag")

    stringFlag, _ := SetString("s", "string", "A string flag", false, "default")
    
    require.NotNil(t, stringFlag, "SetString should return a non-nil pointer")
    require.Equal(t, "", *stringFlag, "Default value of a string flag should be empty string")

    flag, exists := flags.getByLongName["--string"]
    require.True(t, exists, "Flag should be stored in the flags map")
    require.Equal(t, "-s", flag.ShortFormName, "Short form name should match")
    require.Equal(t, "--string", flag.LongFormName, "Long form name should match")
    require.Equal(t, "A string flag", flag.Description, "Description should match")
    require.Equal(t, false, flag.Required, "Required status should match")

    defaultValue, exists := defaultValues[stringFlag]
    require.True(t, exists, "Default value should be stored in the defaultValues map")
    require.Equal(t, "default", defaultValue, "Default value should match")

    setMyOwnArgs("--string", "testing")
    Parse()

    require.Equal(t, "testing", *stringFlag, "String flag should be set to default value after parsing")
}

// Remove the stupid default test flags, surely this is best practice
func setMyOwnArgs(args ...string) {
    os.Args = append([]string{os.Args[0]}, args...)
}

