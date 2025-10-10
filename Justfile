set shell := ["bash", "-cu"]

install: test build
    @[ -f /usr/local/bin/yfile ] && echo "Removing existing installation..." && sudo rm -f /usr/local/bin/yfile || true
    @echo "Installing yfile to /usr/local/bin..."
    @sudo install -m 755 build/yfile /usr/local/bin/yfile
    @echo "yfile installed to /usr/local/bin/yfile"

uninstall:
    @sudo rm -f /usr/local/bin/yfile
    @echo "yfile uninstalled from /usr/local/bin/yfile"

build: check-deps clear-build-cache compile-yara
    @go build -o build/yfile ./cmd/yfile

compile-yara:
    @yarac ./rules/index.yar ./pkg/scanning/matcher/ruleset.compiled 2> /dev/null || (echo "Failed to compile YARA rules!" && exit 1)

clear-build-cache:
    @go clean -cache -modcache -r -i

test: check-deps clear-build-cache compile-yara lint
    @go test -v ./...

lint:
    @golangci-lint run ./...
    @editorconfig-checker

check-deps:
    @command -v yara > /dev/null || (echo "Command 'yara' not found: please install yara and libyara-dev from your package manager!" && exit 1)
    @command -v yarac > /dev/null || (echo "Command 'yarac' not found: please install yara and libyara-dev from your package manager!\n" && exit 1)
    @command -v golangci-lint > /dev/null || (echo "Command 'golangci-lint' not found: please run 'just dev-setup' to install it!" && exit 1)
    @command -v editorconfig-checker > /dev/null || (echo "Command 'editorconfig-checker' not found: please run 'just dev-setup' to install it!" && exit 1)

dev-setup:
    @go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    @go install github.com/editorconfig-checker/editorconfig-checker/cmd/editorconfig-checker@latest

get url:
    @go get "{{url}}"
    @go mod tidy

release: build
    @mkdir -p ./build/release
    @cp ./build/yfile ./yfile
    @tar -cvzf build/release/yfile-linux-amd64.tar.gz ./yfile
    @rm ./yfile
