package utils

import (
    "path/filepath"
    "fmt"
    "os"
)

func IsDir(filePath string) bool {
    fileInfo, err := getFileInfo(filePath)
    if err != nil {
        return false
    }

    return fileInfo.IsDir()
}

func IsFile(filePath string) bool {
    fileInfo, err := getFileInfo(filePath)
    if err != nil {
        return false
    }

    return fileInfo.Mode().IsRegular()
}

func ExtractFilesFromDir(dirPath string, recurse *bool) ([]string, error) {
    fileInfo, err := getFileInfo(dirPath)
    if err != nil {
        return []string{}, err
    }

    if !fileInfo.IsDir() {
        return []string{}, fmt.Errorf("path is not a directory: %s", dirPath)
    }
    dirEntries, err := os.ReadDir(dirPath)
    if err != nil {
        return []string{}, err
    }
    files := make([]string, 0)
    for _, entry := range dirEntries {
        fullPath := filepath.Join(dirPath, entry.Name())
        if !entry.IsDir() {
            files = append(files, fullPath)
        } else {
            if !*recurse {
                continue
            }

            subDirFiles, err := ExtractFilesFromDir(fullPath, recurse)
            if err != nil {
                fmt.Fprintf(os.Stderr, "[Warning] Could not read sub-directory %s: %v\n", entry.Name(), err)
                continue
            }
            files = append(files, subDirFiles...)
        }
    }

    return files, nil
}

func VerifyFilePath(filePath string) bool {
    if _, err := getFileInfo(filePath); err != nil {
        return false
    }

    return true
}

func getFileInfo(filePath string) (os.FileInfo, error) {
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, err
        }
        fmt.Fprintf(os.Stderr, "[Error] Could not retrieve file information for %s: %v\n", filePath, err)
        return nil, err
    }
    return fileInfo, nil
}

