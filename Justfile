set shell := ["bash", "-cu"]

build: compile-yara
    go build -o build/yfile ./cmd/yfile

compile-yara:
    cd ./rules
    yarac ./index_all.yar ../pkg/scanning/ruleset.compiled
    cd ..

test: lint
    go test -v ./...

lint:
    golangci-lint run ./cmd/yfile
    editorconfig-checker

dev-setup:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/editorconfig-checker/editorconfig-checker/cmd/editorconfig-checker@latest

push msg:
    git add .
    git commit -m "{{msg}}"
    git push

pull:
    branch=$(git rev-parse --abbrev-ref HEAD)
    git pull origin "$branch"

get url:
    go get "{{url}}"
    go mod tidy
