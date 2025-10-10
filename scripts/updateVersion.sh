#!/usr/bin/env bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <new_version>"
    exit 1
elif [[ ! $1 =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version must be in the format X.Y.Z where X, Y, and Z are integers."
    exit 1
fi

VERSION_FILE=./VERSION
echo "$1" > "$VERSION_FILE"

UTILS_FILE=./pkg/utils/yfileinfo.go
sed -i -E "s|const Version = \"[0-9]+\.[0-9]+\.[0-9]+\"|const Version = \"$1\"|" "$UTILS_FILE"
