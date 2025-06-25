package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobwas/glob"
	"github.com/ocuroot/ifdiff"
)

func main() {
	flagBaseRef := flag.String("base", "HEAD", "The reference to diff against.")
	flagCurrentRef := flag.String("current", "", "The reference to diff. If empty, modified and untracked files are included.")
	noList := flag.Bool("nolist", false, "If set, do not list the files that match the glob.")
	zeroExit := flag.Bool("zero", false, "If set, exit with code 0 even if no files match.")
	flagHelp := flag.Bool("help", false, "Show help message.")
	flag.Parse()

	if *flagHelp {
		usage()
		os.Exit(0)
	}

	// Read globs and follow-on command as provided
	var globs []string
	var cmd []string
	for i, arg := range flag.Args() {
		if arg == "--" {
			cmd = flag.Args()[i+1:]
			break
		}
		globs = append(globs, arg)
	}

	if len(globs) == 0 {
		fmt.Fprintln(os.Stderr, "Error: must specify at least one glob.")
		usage()
		os.Exit(2)
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(2)
	}

	// Identify changed files
	changes, err := ifdiff.ChangedFiles(dir, *flagBaseRef, *flagCurrentRef)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(2)
	}

	// Match changed files to all globs
	matched, err := match(globs, changes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(2)
	}

	if !*noList {
		for _, file := range matched {
			fmt.Println(file)
		}
	}

	if len(matched) == 0 {
		fmt.Fprintln(os.Stderr, "No changes matched.")
		if len(cmd) > 0 || *zeroExit {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if len(cmd) > 0 {
		var args []string
		if len(cmd) > 1 {
			args = cmd[1:]
		}
		c := exec.Command(cmd[0], args...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		execErr := c.Run()
		if execErr != nil {
			fmt.Fprintln(os.Stderr, "Error executing command:", execErr)
			os.Exit(2)
		}
		os.Exit(0)
	}
}

func match(patterns []string, files []string) ([]string, error) {
	var matches []string
	for _, pattern := range patterns {
		g, err := glob.Compile(pattern, filepath.Separator)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if g.Match(file) {
				matches = append(matches, file)
			}
		}
	}
	return matches, nil
}

func usage() {
	fmt.Println("ifdiff is a tool for running commands conditionally based on git changes.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("ifdiff [options] <glob>... [--] <command>...")
	flag.PrintDefaults()

	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("ifdiff --base=HEAD~1 --current=HEAD '**/*.go'")
	fmt.Println("\tLists .go files that have changed in the most recent commit. If no files are listed, the command exits with code 1.")
	fmt.Println("ifdiff --base=HEAD~1 --current=HEAD '**/*.go' -- go test ./...")
	fmt.Println("\tRuns 'go test ./...' if any .go files have changed in the most recent commit.")
	fmt.Println("ifdiff '**/*.go' -- go fmt")
	fmt.Println("\tRuns 'go fmt' if any .go files have been modified but not committed.")
	fmt.Println()
}
