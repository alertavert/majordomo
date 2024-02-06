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

const LocationEnv = "MAJORDOMO_CONFIG"

var DefaultConfigLocation = os.Getenv("HOME") + "/.majordomo/config.yaml"

type Project struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Location    string `yaml:"location" json:"location"`
}

// String function makes the Project type a valid fmt.Stringer
func (p Project) String() string {
	return fmt.Sprintf("Project [Name: %s, Description: %s, Location: %s]", p.Name, p.Description, p.Location)
}

type Config struct {
	LoadedFrom        string `yaml:"-"`
	OpenAIApiKey      string `yaml:"api_key"`
	ScenariosLocation string `yaml:"scenarios"`
	CodeSnippetsDir   string `yaml:"code_snippets"`
	Model             string `yaml:"model"`

	ActiveProject string    `yaml:"active_project"`
	Projects      []Project `yaml:"projects"`
}

// Save writes the Config to a YAML file at the given filePath.
// If filePath is empty, it will write to the location from which the
// Config was loaded.
// FIXME: we need to protect multiple writers to the same file using a mux.
func (c *Config) Save(filepath string) error {
	data, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}
	if filepath == "" {
		filepath = c.LoadedFrom
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	return nil
}

// LoadConfig reads the YAML file at the given filepath and returns a Config
// struct.
// If filepath is empty, it will read from the default location, unless the
// MAJORDOMO_CONFIG environment variable is set, in which case it will read
// from that location.
// FIXME: we need to protect races with a writer to the same file using a mux.
func LoadConfig(filepath string) (*Config, error) {
	var c Config

	if filepath == "" {
		filepath = os.Getenv(LocationEnv)
		if filepath == "" {
			filepath = DefaultConfigLocation
		}
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling yaml: %w", err)
	}

	if len(c.Projects) == 0 {
		// TODO: we should have a default project.
		return nil, fmt.Errorf("no projects configured")
	}

	c.LoadedFrom = filepath
	// Converts relative paths in the test_config.yaml to absolute paths
	// by pre-pending the path to the config file.
	baseDir := path.Dir(filepath)
	if !path.IsAbs(c.ScenariosLocation) {
		c.ScenariosLocation = path.Join(baseDir, c.ScenariosLocation)
	}
	if !path.IsAbs(c.CodeSnippetsDir) {
		c.CodeSnippetsDir = path.Join(baseDir, c.CodeSnippetsDir)
	}
	if c.ActiveProject == "" && len(c.Projects) > 0 {
		c.ActiveProject = c.Projects[0].Name
	}
	return &c, nil
}

func (c *Config) GetProject(name string) *Project {
	for _, p := range c.Projects {
		if p.Name == name {
			return &p
		}
	}
	return nil
}

func (c *Config) GetActiveProject() *Project {
	return c.GetProject(c.ActiveProject)
}
