# yfile

Reimplementation of the "file" command with yara rules to scan for malware signatures and other patterns.


## Requirements

### Development
- `just`
- `yara`
- `libyara-dev`

## Installation

### Releases

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
