/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type Config struct {
	OpenAIApiKey string `yaml:"api_key"`
	ScenariosLocation string `yaml:"scenarios"`
}

func LoadConfig() (Config, error) {
	var c Config

	filePath := os.Getenv("MAJORDOMO_CONFIG")
	if filePath == "" {
		filePath = os.Getenv("HOME") + "/.openai/config.yaml"
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c, fmt.Errorf("error reading config file: %w", err)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return c, fmt.Errorf("error unmarshaling yaml: %w", err)
	}
	return c, nil
}
