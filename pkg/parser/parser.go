/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package parser

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
)
// ParseBotResponse parses the bot response and extracts code snippets
func ParseBotResponse(botSays string) error {
	snippetRegex := regexp.MustCompile(`'''([\w/.]+?)\n([\s\S]+?)'''`)
	matches := snippetRegex.FindAllStringSubmatch(botSays, -1)

	if len(matches) == 0 {
		return errors.New("no valid code snippets found")
	}

	for _, match := range matches {
		filePath := match[1]
		content := match[2]

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}
		}

		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
