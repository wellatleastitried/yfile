# Documentation
For information on how to contribute to or install `yfile`, see the [README](../README.md).

## Use cases
`yfile` serves as a way to quickly scan files to gather intel on their true file type, in addition to detecting whether they contain common malware signatures.
This program can also be used as an API through other apps due to its handling of exit codes:
- 0: Successful execution - file(s) are clean
- 1: Successful execution - file(s) contain malware signatures
- 2: Error

## Example Output
For multiple examples of output, please see the [examples doc](examples.md).
