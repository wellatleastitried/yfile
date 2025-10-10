# Example Output

### Help Menu
```md
\> yfile --help
Usage: ./build/yfile [options] <file1> <file2> ...

Options:
  -h   --help               Show this help message and exit
  -v   --verbose            Enable verbose output
       --version            Show version information and exit
  -f   --file-args          Arguments to pass through to the `file` command (e.g. '-b -i')
  -fh  --file-help          Show help for the `file` command and exit
```

### Version Info
```md
\> yfile --version
yfile version 1.0.0
```

### File Analysis (clean file)
```md
\> yfile example.md
example.md: ASCII text
File does not match common malware signatures.
```

### File Analysis (malicious file)
```md
\> yfile test/lua_malware_sig.lua
test/lua_malware_sig.lua: ASCII text
1 YARA matches:
  - Rule: LuaBot (Namespace: default)
```

### Verbose File Analysis (malicious file)
```md
\> yfile --verbose test/lua_malware_sig.lua
test/testsignatures/lua_malware_sig.lua: ASCII text
1 YARA matches:
  - Rule: LuaBot (Namespace: default)
    Tags: [MALW]
```

### Verbose File Analysis (malicious file with multiple matches)
```md
\> yfile --verbose test/multiple.exe.txt
test/testsignatures/multiple.exe.txt: ASCII text
3 YARA matches:
  - Rule: Contains_UserForm_Object (Namespace: default)
  - Rule: LuaBot (Namespace: default)
    Tags: [MALW]
  - Rule: php_anuna (Namespace: default)
```

### File Analysis with custom `file` arguments
```md
\> yfile --file-args '-b -i' README.md
TODO
```

