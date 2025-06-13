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

func TestIntegration(t *testing.T) {
	// Read API key from .env.test.local file
	var err error
	ApiKey, err = readAPIKeyFromFile(EnvFilePath)
	if err != nil || ApiKey == "" {
		t.Fatalf("Failed to read OPENAI_API_KEY from %s: %v", EnvFilePath, err)
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
