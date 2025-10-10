package utils

const (
    ExitOk = 0
    ExitInfected = 1
    ExitError = 2
)

func GetExitcodeFromMatches(matches int) int {
    if matches > 0 {
        return ExitInfected
    }
    return ExitOk
}
