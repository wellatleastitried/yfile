package main

import (
    "os"
    "strings"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestGetFileInfo(t *testing.T) {
    // Test with a valid file path
    tmpFile, err := os.CreateTemp("", "testfile")
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }
    defer os.Remove(tmpFile.Name())

    filePath := tmpFile.Name()
    info, err := getFileInfo(&filePath)
    if err != nil {
        t.Fatalf("getFileInfo() failed for valid path: %v", err)
    }
    if info.Name() != strings.TrimPrefix(filePath, os.TempDir()+string(os.PathSeparator)) {
        t.Errorf("Expected file name %s, got %s", strings.TrimPrefix(filePath, os.TempDir()+string(os.PathSeparator)), info.Name())
    }

    // Test with a non-existent file path
    invalidPath := "/non/existent/file/path"
    _, err = getFileInfo(&invalidPath)
    if err == nil {
        t.Fatal("Expected error for non-existent path, got nil")
    }

    // Test with an empty file path
    emptyPath := ""
    _, err = getFileInfo(&emptyPath)
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
    validResult := verifyFilePath(&filePath)
    require.Equal(t, true, validResult)

    // Test with a non-existent file path
    invalidPath := "/non/existent/file/path"
    invalidResult := verifyFilePath(&invalidPath)
    require.Equal(t, false, invalidResult)

    // Test with an empty file path
    emptyPath := ""
    emptyResult := verifyFilePath(&emptyPath)
    require.Equal(t, false, emptyResult)
}

