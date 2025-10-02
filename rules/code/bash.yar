rule bash_source_file {
    meta:
        author = "wellatleastitried"
        description = "Determines whether the file is bash source code"
        last_modified = "2025-10-01"
    strings:
        $shebang = "#!/bin/bash"
        $moderShebang = "#!/usr/bin/env bash"
    condition:
        $shebang or $modernShebang
}
