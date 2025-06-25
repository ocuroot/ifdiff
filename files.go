package ifdiff

import (
	"sort"
	"strings"

	"github.com/ocuroot/gittools"
)

type ChangedFilesOptions struct {
	Base    string
	Current string
}

// ChangedFiles returns a list of files that have change between the base and current commit.
// The list is sorted alphabetically.
// repoPath is the path to the git repository.
// base is the earliest commit to diff against.
// current is the latest commit to diff against.
// If current is empty, uncomitted files, including untracked files, are included.
func ChangedFiles(repoPath, base, current string) ([]string, error) {
	var files []string

	repo, err := gittools.Open(repoPath)
	if err != nil {
		return nil, err
	}

	if current == "" {
		untracked, err := repo.LsFiles(gittools.LsFilesOptions{
			Others:          true,
			ExcludeStandard: true,
		})
		if err != nil {
			return nil, err
		}
		untracked = strings.TrimSpace(untracked)
		if untracked != "" {
			files = append(files, strings.Split(untracked, "\n")...)
		}
	}

	var commits = []string{}
	commits = append(commits, base)
	if current != "" {
		commits = append(commits, current)
	}

	out, err := repo.Diff(gittools.DiffOptions{
		NameOnly: true,
	}, commits...)
	if err != nil {
		return nil, err
	}
	out = strings.TrimSpace(out)
	if out != "" {
		files = append(files, strings.Split(out, "\n")...)
	}
	sort.Strings(files)
	return files, nil
}
