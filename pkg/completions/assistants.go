/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package completions

import (
	"gopkg.in/yaml.v3"
	"os"
)


// Assistants is a struct that contains the data necessary to instantiate Assistants.
type Assistants struct {
	Common       string            `yaml:"common"`
	Instructions map[string]string `yaml:"instructions"`
}

// GetInstructions is a method that returns the instructions for a given assistant.
func (s *Assistants) GetInstructions(name string) string {
	if s.Instructions == nil {
		return ""
	}
	return s.Instructions[name]
}

// Names returns the names of all the configured assistants.
//
// This may not necessarily accurately reflect those configured in OpenAI:
// use the `/assistants` API to get the list of available assistants.
func (s *Assistants) Names() []string {
	var names []string
	for name := range s.Instructions {
		names = append(names, name)
	}
	return names
}

// ReadInstructions reads a YAML file which contains instructions to create the
// OpenAI Assistants.
// `location` is a string indicating the path to the YAML file.
//  It returns an error if the file cannot be read or if the YAML cannot be parsed.
func ReadInstructions(location string) (*Assistants, error) {
	var assistants Assistants
	bytes, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &assistants)
	if err != nil {
		return nil, err
	}
	return &assistants, nil
}
