/*
 * Copyright (c) 2024 AlertAvert.com. All rights reserved.
 */

package completions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"

	. "github.com/alertavert/gpt4-go/pkg/completions"
)

var _ = Describe("Integration Tests: When querying OpenAI", func() {

	BeforeEach(func() {
		if apiKey == "" {
			Skip("No OPENAI_API_KEY found")
		}

		// We want to see the logs for integration tests, at least until they become stable
		// TODO: increase the log level to ErrorLevel once the tests are stable
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	})
	AfterEach(func() {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	})

	Context("with a valid API key", func() {
		// Enable logging for the test.
		It("can create a new thread", func() {
			tid := activeBot.CreateNewThread()
			Expect(tid).NotTo(BeEmpty())
		})
		It("should return a response for a valid prompt", func() {
			prompt := "Please update this code:\n'''sample/main.go\n" +
				"'''to also print the current date."
			request := PromptRequest{
				Assistant: "go_developer",
				ThreadId:  "",
				Prompt:    prompt,
			}
			Eventually(func(g Gomega) {
				// TODO: run this in a goroutine and check the response
				// in the main thread.

				response, err := activeBot.QueryBot(&request)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(response).NotTo(BeEmpty())
				// TODO: check that the response contains the expected code.
				// TODO: check that the snippet was saved to the correct location.
				// TODO: once background processing is enabled, change the timeout
			}, "2s", "2s").Should(Succeed())
		})
	})
})
