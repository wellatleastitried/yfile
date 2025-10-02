rule java_source_file {
    meta:
        author = "wellatleastitried"
        description = "Determines if the file is Java source code"
        last_modified = "2025-10-01"
    strings:
        $mainSig = "public static void main(String[]"
    condition:
        $mainSig
}
