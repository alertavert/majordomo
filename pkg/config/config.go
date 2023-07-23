/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var DefaultConfigLocation = os.Getenv("HOME") + "/.majordomo/config.yaml"

type Config struct {
	OpenAIApiKey string `yaml:"api_key"`
	ScenariosLocation string `yaml:"scenarios"`
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
	return c, nil
}
