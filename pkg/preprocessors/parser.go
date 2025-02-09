/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package preprocessors

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	ErrorNoCodeSnippetsFound = "no code found for %s: %v"

	// TODO: replace these Regex patterns with a more robust solution, like
	// 	inline commands (such as LOAD, SAVE, etc.)
	// 	See #19

	FilepathPattern    = `^/?([\w.-]+/?)+$`
	CodeSnippetPattern = `'''([\w/.]+/?)\n([\s\S]+?)'''`
	PromptCodePattern  = `'''([\w/.]+/?)\n'''`
)

// SourceCodeMap is a map of file paths to their contents
type SourceCodeMap = map[string]string

// Parser parses code snippets into a prompt or from a bot response
type Parser struct {
	CodeMap SourceCodeMap
}

// A CodeStoreHandler interface abstracts the storage layer for the code
// snippets via a SourceCodeMap.
type CodeStoreHandler interface {
	// GetSourceCode fills in the code snippets, given their file paths
	GetSourceCode(codemap *SourceCodeMap) error

	// PutSourceCode will store the code snippets, based on their file paths
	PutSourceCode(codemap SourceCodeMap) error
}

var validPathPattern *regexp.Regexp
var snippetRegex *regexp.Regexp
var promptRegex *regexp.Regexp

func init() {
	validPathPattern = regexp.MustCompile(FilepathPattern)
	snippetRegex = regexp.MustCompile(CodeSnippetPattern)
	promptRegex = regexp.MustCompile(PromptCodePattern)
}

func IsValidFilePath(path string) bool {
	return validPathPattern.MatchString(path)
}

// ParseBotResponse parses a prompt or bot response and extracts code snippets
// with their respective file paths.
func (p *Parser) ParseBotResponse(botSays string) error {
	matches := snippetRegex.FindAllStringSubmatch(botSays, -1)
	if len(matches) == 0 {
		return nil
	}

	if p.CodeMap == nil {
		p.CodeMap = make(map[string]string)
	}
	for _, match := range matches {
		filePath := match[1]
		if IsValidFilePath(filePath) {
			content := match[2]
			p.CodeMap[filePath] = content
		} else {
			return errors.New(fmt.Sprintf("invalid file path: %s", filePath))
		}
	}
	return nil
}

// ParsePrompt finds all the code snippets in the prompt and extracts their paths
// from the prompt to prepare the CodeMap to be populated by a CodeStoreHandler.
func (p *Parser) ParsePrompt(prompt string) {
	matches := promptRegex.FindAllStringSubmatch(prompt, -1)
	for _, match := range matches {
		p.CodeMap[match[1]] = ""
	}
}

// FillPrompt fills in the code snippets in the prompt, given their file paths.
func (p *Parser) FillPrompt(prompt string) (string, error) {
	matches := promptRegex.FindAllStringSubmatch(prompt, -1)

	for _, match := range matches {
		filePath := match[1]
		content, found := p.CodeMap[filePath]
		if content == "" || !found {
			return "", errors.New(fmt.Sprintf(ErrorNoCodeSnippetsFound, filePath,
				"no entry in map"))
		}
		replacementRegex := regexp.MustCompile(fmt.Sprintf(`'''%s\n'''`, filePath))
		prompt = replacementRegex.ReplaceAllLiteralString(prompt,
			fmt.Sprintf("'''%s\n%s'''", filePath, content))
	}
	return prompt, nil
}
