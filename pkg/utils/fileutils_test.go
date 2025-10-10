package utils

import (
    "os"
    "strings"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestIsFile(t *testing.T) {
    // Test with a valid file path
    tmpFile, err := os.CreateTemp("", "testfile")
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }
    defer os.Remove(tmpFile.Name())

    filePath := tmpFile.Name()
    isFile := IsFile(filePath)
    require.Equal(t, true, isFile)

    // Test with a directory path
    dirPath := os.TempDir()
    isFile = IsFile(dirPath)
    require.Equal(t, false, isFile)

    // Test with a non-existent file path
    invalidPath := "/non/existent/file/path"
    isFile = IsFile(invalidPath)
    require.Equal(t, false, isFile)

    // Test with an empty file path
    emptyPath := ""
    isFile = IsFile(emptyPath)
    require.Equal(t, false, isFile)
}

func TestIsDir(t *testing.T) {
    // Test with a valid directory path
    dirPath := os.TempDir()
    isDir := IsDir(dirPath)
    require.Equal(t, true, isDir)

    // Test with a file path
    tmpFile, err := os.CreateTemp("", "testfile")
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }
    defer os.Remove(tmpFile.Name())

    filePath := tmpFile.Name()
    isDir = IsDir(filePath)
    require.Equal(t, false, isDir)

    // Test with a non-existent directory path
    invalidPath := "/non/existent/directory/path"
    isDir = IsDir(invalidPath)
    require.Equal(t, false, isDir)

    // Test with an empty directory path
    emptyPath := ""
    isDir = IsDir(emptyPath)
    require.Equal(t, false, isDir)
}

func TestExtractFilesFromDir(t *testing.T) {
    // Test with a valid directory path
    dirPath, err := os.MkdirTemp(os.TempDir(), "testdir")
    if err != nil {
        t.Fatalf("Failed to create temp directory: %v", err)
    }
    defer os.RemoveAll(dirPath)

    files, err := ExtractFilesFromDir(dirPath)
    if err != nil {
        t.Fatalf("ExtractFilesFromDir() failed for valid directory: %v", err)
    }
    require.GreaterOrEqual(t, len(files), 0)

    // Test with a directory containing files and sub-directories
    _, err = os.CreateTemp(dirPath, "file1")
    if err != nil {
        t.Fatalf("Failed to create temp file in test directory: %v", err)
    }

    subDir := dirPath + "/subdir"

    err = os.Mkdir(subDir, 0755)
    if err != nil {
        t.Fatalf("Failed to create sub-directory: %v", err)
    }
    _, err = os.CreateTemp(subDir, "file2")
    if err != nil {
        t.Fatalf("Failed to create temp file in sub-directory: %v", err)
    }

    files, err = ExtractFilesFromDir(dirPath)
    if err != nil {
        t.Fatalf("ExtractFilesFromDir() failed for directory with files and sub-directories: %v", err)
    }

    require.Equal(t, 2, len(files))
    require.Contains(t, files[0], "file1")
    require.Contains(t, files[1], "file2")

    // Test with a file path instead of a directory
    tmpFile, err := os.CreateTemp("", "testfile")
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }
    defer os.Remove(tmpFile.Name())

    filePath := tmpFile.Name()
    _, err = ExtractFilesFromDir(filePath)
    if err == nil {
        t.Fatal("Expected error for file path, got nil")
    }

    // Test with a non-existent directory path
    invalidPath := "/non/existent/directory/path"
    _, err = ExtractFilesFromDir(invalidPath)
    if err == nil {
        t.Fatal("Expected error for non-existent directory, got nil")
    }

    // Test with an empty directory path
    emptyPath := ""
    _, err = ExtractFilesFromDir(emptyPath)
    if err == nil {
        t.Fatal("Expected error for empty directory path, got nil")
    }
}

func TestGetFileInfo(t *testing.T) {
    // Test with a valid file path
    tmpFile, err := os.CreateTemp("", "testfile")
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }
    defer os.Remove(tmpFile.Name())

    filePath := tmpFile.Name()
    info, err := getFileInfo(filePath)
    if err != nil {
        t.Fatalf("getFileInfo() failed for valid path: %v", err)
    }
    if info.Name() != strings.TrimPrefix(filePath, os.TempDir()+string(os.PathSeparator)) {
        t.Errorf("Expected file name %s, got %s", strings.TrimPrefix(filePath, os.TempDir()+string(os.PathSeparator)), info.Name())
    }

    // Test with a non-existent file path
    invalidPath := "/non/existent/file/path"
    _, err = getFileInfo(invalidPath)
    if err == nil {
        t.Fatal("Expected error for non-existent path, got nil")
    }

    // Test with an empty file path
    emptyPath := ""
    _, err = getFileInfo(emptyPath)
    if err == nil {
        t.Fatal("Expected error for empty path, got nil")
    }
}

func TestVerifyFilePath(t *testing.T) {
    // Test with a valid file path
    tmpFile, err := os.CreateTemp("", "testfile")
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }
    defer os.Remove(tmpFile.Name())

    filePath := tmpFile.Name()
    validResult := VerifyFilePath(filePath)
    require.Equal(t, true, validResult)

    // Test with a non-existent file path
    invalidPath := "/non/existent/file/path"
    invalidResult := VerifyFilePath(invalidPath)
    require.Equal(t, false, invalidResult)

    // Test with an empty file path
    emptyPath := ""
    emptyResult := VerifyFilePath(emptyPath)
    require.Equal(t, false, emptyResult)
}

