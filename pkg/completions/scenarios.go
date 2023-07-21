/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package completions

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// cached scenarios
var scenarios Scenarios

// Scenarios is a struct that contains the data from the scenarios.yaml file.
type Scenarios struct {
	Common    string            `yaml:"common,omitempty"`
	Scenarios map[string]string `yaml:"scenarios,omitempty"`
}

func (s *Scenarios) GetScenario(name string) string {
	if s.Scenarios == nil {
		return ""
	}
	return s.Scenarios[name]
}

func (s *Scenarios) GetCommon() string {
	return s.Common
}

func (s *Scenarios) GetScenarioNames() []string {
	var names []string
	for name := range s.Scenarios {
		names = append(names, name)
	}
	return names
}

// ReadScenarios is a function that reads a YAML file which contains scenarios and
//  caches the data in the `scenarios` package variable.
// `location` is a string indicating the path to the YAML file.
//  It returns an error if the file cannot be read or if the YAML cannot be unmarshaled.
//
// Use `GetScenarios()` to retrieve the cached scenarios.
func ReadScenarios(location string) error {
	bytes, err := os.ReadFile(location)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, &scenarios)
	if err != nil {
		return err
	}

	return nil
}

// GetScenarios is a function that returns a pointer to the cached scenarios.
func GetScenarios() *Scenarios {
	if scenarios.Scenarios == nil {
		err := ReadScenarios(os.Getenv("MAJORDOMO_SCENARIOS"))
		if err != nil {
			fmt.Println("Error reading scenarios file: ", err)
			return nil
		}
	}
	return &scenarios
}
