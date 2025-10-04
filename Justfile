set shell := ["bash", "-cu"]

build:
    go build -o build/yfile ./cmd/yfile

test: lint
    go test -v ./...

lint:
    golangci-lint run ./cmd/yfile
    editorconfig-checker

dev-setup:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@1.64.8
    go install github.com/editorconfig-checker/editorconfig-checker/cmd/editorconfig-checker@latest

push msg:
    git add .
    git commit -m "{{msg}}"
    git push

pull:
    branch=$(git rev-parse --abbrev-ref HEAD)
    git pull origin "$branch"

