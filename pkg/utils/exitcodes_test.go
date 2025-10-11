package utils

import (
    "testing"

    "github.com/stretchr/testify/require"
)

func TestGetExitcodeFromMatches(t *testing.T) {
    noMatches := 0
    oneMatch := 1
    manyMatches := 48

    require.Equal(t, ExitOk, GetExitcodeFromMatches(noMatches))
    require.Equal(t, ExitInfected, GetExitcodeFromMatches(oneMatch))
    require.Equal(t, ExitInfected, GetExitcodeFromMatches(manyMatches))
}
