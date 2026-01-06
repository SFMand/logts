# logts (Log Timestamped)
A CLI tool written in Go, for archiving and compressing directories. Built using `cobra-cli`, this tool creates timestamped `.tar.gz` archives while preserving directory structure.

This project was created to practice Go and its concurrency features (`goroutines`, `channels`, `waitgroups`) and to learn how to build Command Line Interfaces.

## Main Topics Covered
- **Concurrency:** Utilized a worker pool pattern to parallelize file opening and header preparation. Also, managed communication between the workers (Processing the files) and the consumer (Tar writer) using unbuffered channels.
- **Memory Safety:** Implemented streaming compression (Gzip -> Tar) to handle large directories, to avoid loading large files into memory.
- **Cross-Platform Compatibility:** Utilized Go's `filepath` functions to handle the differences between file paths in Windows `\` and Linux `/`.

## Prerequistes
- **Go:** Version 1.21 or higher

## Installation
### Options (Pick One)
1. Install via Go: 
      
        go install github.com/SFMand/logts@latest

2. Build the binary file from source:

        git clone https://github.com/SFMand/logts.git
        cd logts
        go build -o main.go

### Usage
The tool accepts a source directory and a destination path, source is required, but destination can be omitted, which would produce the archive in cwd.

        logts --from [source_directory] --to [destination_path]

**Example**

        logts --from /var/log/example --to ~/backups
**Output**

A compressed folder in `~/backups` with name `logts_example-timestamp.tar.gz` is generated.
