package utils

import (
    "runtime"
    "strconv"
)

const Version = "0.1.1"

func CPUCountString() string {
    return strconv.Itoa(runtime.NumCPU())
}

// Min 1, default 6, max runtime.NumCPU()
func GetWorkerCount(numWorkers *string) int {
    if numWorkers == nil || *numWorkers == "" {
        return 6
    }

    workers, err := strconv.Atoi(*numWorkers)
    if err != nil || workers < 1 {
        return 1
    }
    maxWorkers := runtime.NumCPU()
    if workers > maxWorkers {
        return maxWorkers
    }

    return workers
}
