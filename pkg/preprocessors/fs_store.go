/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package preprocessors

import (
	"errors"
	"fmt"
	"github.com/alertavert/gpt4-go/pkg/config"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

const (
	ErrorReadingCodeSnippet = "error while reading %s: %v"
)

type ProjectsStoreMap = map[string]*CodeStoreHandler

var cache = make(ProjectsStoreMap)

// FilesystemStore is a CodeStoreHandler that reads and writes code snippets from/to the filesystem
type FilesystemStore struct {
	// SourceCodeDir is the directory where the code snippets are read from
	SourceCodeDir string
	// DestCodeDir is the directory where the code snippets are saved to
	DestCodeDir string
}

func (fp *FilesystemStore) GetSourceCode(codeMap *SourceCodeMap) error {
	for relPath := range *codeMap {
		content, err := os.ReadFile(filepath.Join(fp.SourceCodeDir, relPath))
		if err != nil {
			return errors.New(fmt.Sprintf(ErrorReadingCodeSnippet, relPath, err))
		}
		(*codeMap)[relPath] = string(content)
	}
	return nil
}

func (fp *FilesystemStore) PutSourceCode(codemap SourceCodeMap) error {
	for relPath, content := range codemap {
		absPath := filepath.Join(fp.DestCodeDir, relPath)
		dir := filepath.Dir(absPath)
		// Creates the directory if it doesn't exist
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
			log.Debug().
				Str("path", dir).
				Msg("new directory created to store snippets")
		}
		// Writes to the file in its respective directory
		err := os.WriteFile(absPath, []byte(content), 0644)
		if err != nil {
			return err
		}
		log.Debug().
			Str("path", absPath).
			Str("relative_path", relPath).
			Msg("Code saved to file")
	}
	return nil
}

// NewFilesystemStore creates a new filesystem-based CodeStoreHandler
func NewFilesystemStore(sourceDir, destDir string) CodeStoreHandler {
	return &FilesystemStore{
		SourceCodeDir: sourceDir,
		DestCodeDir:   destDir,
	}
}

// GetCodeStoreHandler returns a CodeStoreHandler for the given project
// Creating a new one if necessary.
func GetCodeStoreHandler(project *config.Project) *CodeStoreHandler {
	if cache[project.Name] == nil {
		store := NewFilesystemStore(project.Location, project.ResolvedCodeSnippetsDir)
		cache[project.Name] = &store
	}
	return cache[project.Name]
}
