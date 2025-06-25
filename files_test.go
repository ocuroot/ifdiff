package ifdiff

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/ocuroot/gittools"
	"github.com/stretchr/testify/assert"
)

type ValidFileTest struct {
	RepoPath      string
	Base          string
	Current       string
	ExpectedFiles []string
}

func createBasicTests(t *testing.T) []ValidFileTest {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "ifdiff-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	git := gittools.Client{}
	repo, err := git.Init(tmpDir, "main")
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	for i := 0; i < 100; i++ {
		if err := os.WriteFile(filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i)), []byte("Hello world"), 0o666); err != nil {
			t.Fatalf("Failed to create file %d: %v", i, err)
		}
	}
	if err := repo.CommitAll("Initial commit"); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	var expectedFiles []string
	for i := 0; i < 5; i++ {
		expectedFiles = append(expectedFiles, fmt.Sprintf("file%d.txt", i))
		if err := os.WriteFile(filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i)), []byte("Updated"), 0o666); err != nil {
			t.Fatalf("Failed to update file %d: %v", i, err)
		}
	}
	if err := repo.CommitAll("Edit 5 files"); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}
	sort.Strings(expectedFiles)

	var modifiedAndUntracked []string
	for i := 5; i < 10; i++ {
		modifiedAndUntracked = append(modifiedAndUntracked, fmt.Sprintf("file%d.txt", i))
		if err := os.WriteFile(filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i)), []byte("Updated (uncomitted)"), 0o666); err != nil {
			t.Fatalf("Failed to update file %d: %v", i, err)
		}

		modifiedAndUntracked = append(modifiedAndUntracked, fmt.Sprintf("file%d.txt", 100+i))
		if err := os.WriteFile(filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", 100+i)), []byte("New"), 0o666); err != nil {
			t.Fatalf("Failed to update file %d: %v", i, err)
		}
	}
	sort.Strings(modifiedAndUntracked)

	return []ValidFileTest{
		{
			RepoPath:      tmpDir,
			Base:          "HEAD~1",
			Current:       "HEAD",
			ExpectedFiles: expectedFiles,
		},
		{
			RepoPath:      tmpDir,
			Base:          "HEAD",
			Current:       "HEAD",
			ExpectedFiles: []string{},
		},
		{
			RepoPath:      tmpDir,
			Base:          "HEAD",
			Current:       "",
			ExpectedFiles: modifiedAndUntracked,
		},
	}
}

func TestFiles(t *testing.T) {
	tests := createBasicTests(t)

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s..%s", test.Base, test.Current), func(t *testing.T) {
			files, err := ChangedFiles(test.RepoPath, test.Base, test.Current)
			if err != nil {
				t.Fatalf("Failed to get changed files: %v", err)
			}
			assert.Equal(t, test.ExpectedFiles, files)
		})
	}
}
