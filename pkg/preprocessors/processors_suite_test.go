/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package preprocessors_test

import (
	"fmt"
	"github.com/alertavert/gpt4-go/pkg/preprocessors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"testing"
)

const (
	TestdataDir = "../../testdata"
)

func TestParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pre-Processors Test Suite")
}

var _ = BeforeSuite(func() {
	// Silence the logs
	zerolog.SetGlobalLevel(zerolog.Disabled)
})

// SetupTestFiles creates two randomly named directories and copies files from `TestdataDir`
// folder and its sub-folders into the first of them.
// It returns the paths of the two directories, and an error if any.
func SetupTestFiles() (srcDir string, destDir string, err error) {
	srcDir, _ = os.MkdirTemp("", "src")
	destDir, _ = os.MkdirTemp("", "dest")
	// Copying the files from `testdata` to srcDir
	err = filepath.Walk(TestdataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(TestdataDir, path)
		destPath := filepath.Join(srcDir, relPath)
		if !info.IsDir() {
			// Reading the file
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// Writing the file at the destination
			err = os.WriteFile(destPath, data, info.Mode())
			if err != nil {
				return err
			}
		} else {
			err := os.MkdirAll(destPath, info.Mode())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return "", "", err
	}
	Expect(CompareDirs(TestdataDir, srcDir)).To(BeTrue())
	return srcDir, destDir, nil
}

// CleanupTestFiles removes the two temporary directories created by SetupTestFiles
// We ignore errors here, as the directories may not exist if the test failed; and we don't want
// to stop halfway when cleaning up.
func Cleanup(store *preprocessors.FilesystemStore) {
	_ = os.RemoveAll(store.SourceCodeDir)
	_ = os.RemoveAll(store.DestCodeDir)
}

// FillCodemap walks the directory tree rooted at path and fills in the codeMap with the
// relative paths of the files found in the directory.
// If shouldFillContents is true, the contents of the files are also filled in.
func FillCodemap(path string, codeMap preprocessors.SourceCodeMap, shouldFillContents bool) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, err := filepath.Rel(path, filePath)
			if err != nil {
				return err
			}
			var contents []byte
			if shouldFillContents {
				contents, err = os.ReadFile(filePath)
				if err != nil {
					return err
				}
			} else {
				contents = []byte{}
			}
			codeMap[relativePath] = string(contents)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// CompareDirs compares the contents of two directories, and returns true if they are identical.
func CompareDirs(dir1, dir2 string) (bool, error) {
	err := filepath.Walk(dir1, func(path1 string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Generating relative path
			relPath, _ := filepath.Rel(dir1, path1)
			path2 := filepath.Join(dir2, relPath)

			// Checking the file existence in the second directory
			if _, err := os.Stat(path2); os.IsNotExist(err) {
				return fmt.Errorf("file %s in dest dir (%s) does not exist: %v",
					path2, dir2, err)
			}

			// Comparing the file contents
			file1Contents, err := os.ReadFile(path1)
			if err != nil {
				return fmt.Errorf("while reading %s: %s", path1, err)
			}

			file2Contents, err := os.ReadFile(path2)
			if err != nil {
				return fmt.Errorf("while reading %s: %s", path2, err)
			}

			if string(file1Contents) != string(file2Contents) {
				return fmt.Errorf("file contents differ: %s, %s", path1, path2)
			}
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
