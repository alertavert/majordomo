/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

var DefaultConfigLocation = os.Getenv("HOME") + "/.majordomo/config.yaml"

type Config struct {
	OpenAIApiKey      string `yaml:"api_key"`
	ScenariosLocation string `yaml:"scenarios"`
	CodeSnippetsDir   string `yaml:"code_snippets"`
}

func LoadConfig() (Config, error) {
	var c Config

	filePath := os.Getenv("MAJORDOMO_CONFIG")
	if filePath == "" {
		filePath = DefaultConfigLocation
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return c, fmt.Errorf("error reading config file: %w", err)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return c, fmt.Errorf("error unmarshaling yaml: %w", err)
	}

	// Converts relative paths in the test_config.yaml to absolute paths
	// by pre-pending the path to the config file.
	baseDir := path.Dir(filePath)
	if !path.IsAbs(c.ScenariosLocation) {
		c.ScenariosLocation = path.Join(baseDir, c.ScenariosLocation)
	}
	if !path.IsAbs(c.CodeSnippetsDir) {
		c.CodeSnippetsDir = path.Join(baseDir, c.CodeSnippetsDir)
	}
	return c, nil
}
