# ifdiff

[![GoDoc](https://godoc.org/github.com/ocuroot/ifdiff?status.svg)](https://godoc.org/github.com/ocuroot/ifdiff)

A utility to perform actions based on changed files in a Git repository.

## Overview

ifdiff is a tool that helps you run commands conditionally based on which files have changed in your Git repository. It lets you define glob patterns to match against changed files and run subsequent commands only when matches are found.

## Installation

### From Source

```bash
go install github.com/ocuroot/ifdiff/cmd/ifdiff@latest
```

### Pre-built Binaries

You can download pre-built binaries for various platforms from the [GitHub releases page](https://github.com/ocuroot/ifdiff/releases).

## Usage

```
ifdiff [options] <glob>... [--] <command>...
```

### Options

- `-base`: The reference to diff against (default: "HEAD")
- `-current`: The reference to diff (if empty, modified and untracked files are included)
- `-nolist`: If set, do not list the files that match the glob
- `-zero`: If set, exit with code 0 even if no files match
- `-help`: Show help message

### Exit Codes

- `0`: Command executed successfully or no files matched with `-zero` flag
- `1`: No files matched the glob pattern (without `-zero` flag)
- `2`: Error occurred

## Examples

```bash
# List .go files that have changed in the most recent commit
ifdiff --base=HEAD~1 --current=HEAD '**/*.go'

# Run 'go test ./...' if any .go files have changed in the most recent commit
ifdiff --base=HEAD~1 --current=HEAD '**/*.go' -- go test ./...

# Run 'go fmt' if any .go files have been modified but not committed
ifdiff '**/*.go' -- go fmt
```

## Use Cases

- Run tests only for components that have changed
- Execute formatting or linting tools only on modified files
- Trigger build processes for specific components
- Perform custom validation on certain file types when they change
