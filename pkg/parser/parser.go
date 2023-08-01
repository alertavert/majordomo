/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package parser

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
)

const (
	ErrorNoCodeSnippetsFound = "no code found for %s: %v"
	ErrorReadingCodeSnippet  = "error while reading %s: %v"
)

// SourceCode is a map of file paths to their contents
type SourceCode = map[string]string

// ParseBotResponse parses the bot response and extracts code snippets
func ParseBotResponse(botSays string) (SourceCode, error) {
	snippetRegex := regexp.MustCompile(`'''([\w/.]+?)\n([\s\S]+?)'''`)
	matches := snippetRegex.FindAllStringSubmatch(botSays, -1)

	if len(matches) == 0 {
		return nil, errors.New("no valid code snippets found")
	}

	var sourceCode SourceCode = make(map[string]string)
	for _, match := range matches {
		filePath := match[1]
		content := match[2]
		sourceCode[filePath] = content
	}
	return sourceCode, nil
}

// InsertSourceCode inserts the code snippets into prompt text
func InsertSourceCode(prompt string) (string, error) {
	snippetRegex := regexp.MustCompile(`'''([\w/.]+?)\n'''`)
	matches := snippetRegex.FindAllStringSubmatch(prompt, -1)

	for _, match := range matches {
		filePath := match[1]
		f, err := os.Open(filePath)
		if err != nil {
			return "", errors.New(fmt.Sprintf(ErrorNoCodeSnippetsFound, filePath, err))
		}
		content, err := io.ReadAll(f)
		if err != nil {
			return "", errors.New(fmt.Sprintf(ErrorReadingCodeSnippet, filePath, err))
		}
		replacementRegex := regexp.MustCompile(fmt.Sprintf(`'''%s\n'''`, filePath))
		prompt = replacementRegex.ReplaceAllLiteralString(prompt,
			fmt.Sprintf("'''%s\n%s'''", filePath, content))
	}
	return prompt, nil
}
