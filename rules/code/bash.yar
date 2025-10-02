rule bash_source_file {
    meta:
        author = "wellatleastitried"
        description = "Determines whether the file is bash source code"
    strings:
        $shebang = "#!/bin/bash"
        $moderShebang = "#!/usr/bin/env bash"
    condition:
        $shebang or $modernShebang
}
