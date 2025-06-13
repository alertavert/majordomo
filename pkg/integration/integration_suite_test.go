/*
 * Copyright (c) 2024 AlertAvert.com. All rights reserved.
 */

package integration

import (
	"bufio"
	"fmt"
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/alertavert/gpt4-go/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	TestConfigLocation = "../../testdata/test_config.yaml"
	EnvFilePath        = "../../.env.test.local"
)

// This will be instantiated with a valid API key, if found.
var ActiveBot *completions.Majordomo
var ApiKey string

// readAPIKeyFromFile reads the OPENAI_API_KEY from the .env.test.local file
func readAPIKeyFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Err(err).
				Str("file", filePath).
				Msg("Failed to close file")
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			return "", err
		}
		if strings.HasPrefix(line, "OPENAI_API_KEY") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", fmt.Errorf("OPENAI_API_KEY not found in file %s", filePath)
}

// getAPIKey first tries to get the API key from the environment variable,
// and if not found, falls back to reading from the file
func getAPIKey() (string, error) {
	// First try to get the API key from the environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		return apiKey, nil
	}

	// If not found in environment, try to read from file
	return readAPIKeyFromFile(EnvFilePath)
}

func TestIntegration(t *testing.T) {
	// Get API key from environment variable or file
	var err error
	ApiKey, err = getAPIKey()
	if err != nil || ApiKey == "" {
		t.Fatalf("Failed to get OPENAI_API_KEY: %v", err)
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	cfg, err := config.LoadConfig(TestConfigLocation)
	Expect(err).NotTo(HaveOccurred())
	cfg.OpenAIApiKey = ApiKey
	ActiveBot, err = completions.NewMajordomo(cfg)
	Expect(err).NotTo(HaveOccurred())
	Expect(ActiveBot.SetActiveProject("actual")).NotTo(HaveOccurred())
})
