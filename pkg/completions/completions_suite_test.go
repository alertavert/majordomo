/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package completions_test

import (
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/alertavert/gpt4-go/pkg/config"
	"github.com/rs/zerolog"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// This will be instantiated with a valid API key, if found.
var activeBot *completions.Majordomo

func TestCompletions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Completions Suite")
}
var apiKey string

var _ = BeforeSuite(func() {
	// Silence the logs
	zerolog.SetGlobalLevel(zerolog.Disabled)
	apiKey = os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		// This will help debug if the test fails,
		// as it gets emitted even without the -v flag in that case.
		_, _ = GinkgoWriter.Write([]byte("OPENAI_API_KEY found\n"))
		cfg, _ := config.LoadConfig(TestConfigLocation)
		cfg.OpenAIApiKey = apiKey
		activeBot, _ = completions.NewMajordomo(cfg)
		Expect(activeBot.SetActiveProject("actual")).NotTo(HaveOccurred())
	}
})
