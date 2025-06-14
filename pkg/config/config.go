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
const CodeLocationEnv = "MAJORDOMO_CODE"

var DefaultConfigLocation = os.Getenv("HOME") + "/.majordomo/config.yaml"
var DefaultCodeSnippetsLocation = os.Getenv("HOME") + "/.majordomo/code"

type Project struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Location    string `yaml:"location" json:"location"`
	// User-configured location for the code snippets for this project.
	// It is never overwritten by the system.
	CodeSnippets string `yaml:"code_snippets,omitempty" json:"code_snippets,omitempty"`

	// Resolved path for code snippets for the project.
	// This is what the system uses, but is not written to the config file.
	ResolvedCodeSnippetsDir string `yaml:"-" json:"-"`
}

// String function makes the Project type a valid fmt.Stringer
func (p Project) String() string {
	return fmt.Sprintf("Project [Name: %s, Description: %s, Location: %s]", p.Name, p.Description, p.Location)
}

type Config struct {
	// LoadedFrom is the path from which the Config was loaded.
	LoadedFrom string `yaml:"-"`

	// OpenAIApiKey is the API key to use for OpenAI API requests.
	OpenAIApiKey string `yaml:"api_key"`

	// ProjectId is the ID of the project to use for OpenAI API requests.
	ProjectId string `yaml:"project_id"`

	// AssistantsLocation is the path to the YAML file containing the instructions
	// to create the assistants' system prompts.
	// TODO: not supported yet (see #18)
	AssistantsLocation string `yaml:"assistants"`

	// ThreadsLocation is the path to the directory where the conversations are stored.
	ThreadsLocation string `yaml:"threads_location"`

	// CodeSnippetsDir is the name of the directory, inside each respective
	// project's location, where the code snippets are stored.
	CodeSnippetsDir string `yaml:"code_snippets"`

	// Model is the name of the model to use for OpenAI API requests.
	Model string `yaml:"model"`

	// ActiveProject is the name of the project that is currently active and will
	// be used to fetch the files from (and save snippets to).
	ActiveProject string `yaml:"active_project"`

	// Projects is a list of projects that are configured in the system.
	Projects []Project `yaml:"projects"`
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

	// TODO: not having any projects configured should be a valid state.
	if len(c.Projects) == 0 {
		// TODO: we should have a default project.
		return nil, fmt.Errorf("no projects configured")
	}

	c.LoadedFrom = filepath
	// Converts relative paths in the test_config.yaml to absolute paths
	// by pre-pending the path to the config file.
	baseDir := path.Dir(filepath)
	if !path.IsAbs(c.AssistantsLocation) {
		c.AssistantsLocation = path.Join(baseDir, c.AssistantsLocation)
	}
	if c.ActiveProject == "" && len(c.Projects) > 0 {
		c.ActiveProject = c.Projects[0].Name
	}
	// The code snippets, if not configured is read from the environment.
	if c.CodeSnippetsDir == "" {
		c.CodeSnippetsDir = os.Getenv(CodeLocationEnv)
		if c.CodeSnippetsDir == "" {
			c.CodeSnippetsDir = DefaultCodeSnippetsLocation
		}
	}

	// Resolve the actual code snippets directory for each project.
	for i, _ := range c.Projects {
		p := &c.Projects[i]
		// By default, we use the global code snippets directory.
		var cs = c.CodeSnippetsDir
		if p.CodeSnippets != "" {
			// If the project has a code snippets directory configured, we use that.
			cs = p.CodeSnippets
		}
		if !path.IsAbs(cs) {
			// If the path is not absolute, we assume it is relative to the project's location.
			p.ResolvedCodeSnippetsDir = path.Join(p.Location, cs)
		} else {
			p.ResolvedCodeSnippetsDir = cs
		}
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
