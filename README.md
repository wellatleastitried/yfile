# yfile

Reimplementation of the "file" command with yara rules to scan for malware signatures and other patterns.

## Requirements

### Development

#### Packages
- `just`
- `yara`
- `libyara-dev`
- `gcc` / `clang` (Any C compiler)

#### Go Tools
These can be installed by running `just dev-setup` once the project has been cloned.
- `golangci-lint`
- `editorconfig-checker`

## Installation

### Releases
The Binary can be downloaded from the [releases page](https://github.com/wellatleastitried/yfile/releases) OR by running:
```bash
go install github.com/wellatleastitried/yfile/cmd/yfile@latest
```

### Manual
1. Clone the repository with submodules:
`git clone --recurse-submodules https://github.com/wellatleastitried/yfile.git`

2. Run `just`
```bash
just build`
```

3. The binary can be ran from:
```bash
./build/yfile
```
or installed to `/usr/local/bin` with:
```bash
sudo just install
```
