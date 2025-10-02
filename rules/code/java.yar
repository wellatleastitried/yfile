rule java_source_file {
    meta:
        author = "wellatleastitried"
        description = "Determines if the file is Java source code"
    strings:
        $mainSig = "public static void main(String[]"
    condition:
        $mainSig
}
