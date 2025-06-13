/*
 * Copyright (c) 2024 AlertAvert.com. All rights reserved.
 */

package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alertavert/gpt4-go/pkg/completions"
)

var _ = Describe("Integration Tests: When querying OpenAI", func() {
	Context("with a valid API key", func() {
		// Enable logging for the test.
		It("can create a new thread", func() {
			tid := ActiveBot.CreateNewThread("test-project", "go_developer", "test-thread")
			Expect(tid).NotTo(BeEmpty())
		})

		It("should return a response for a valid prompt", func() {
			prompt := "Please update this code:\n'''sample/main.go\n" +
				"'''to also print the current date."
			request := completions.PromptRequest{
				Assistant:  "go_developer",
				ThreadId:   "",
				ThreadName: "test-thread",
				Prompt:     prompt,
			}
			Eventually(func(g Gomega) {
				// TODO: run this in a goroutine and check the response
				// in the main thread.

				response, err := ActiveBot.QueryBot(&request)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(response).NotTo(BeEmpty())
				// TODO: check that the response contains the expected code.
				// TODO: check that the snippet was saved to the correct location.
				// TODO: once background processing is enabled, change the timeout
			}, "2s", "2s").Should(Succeed())
		})

		It("should create a new thread even if the project had never been seen before", func() {
			tid := ActiveBot.CreateNewThread("non-existent-project", "not an assistant", "non-existent-thread")
			Expect(tid).NotTo(BeEmpty())
		})

		It("requires a valid project and assistant to create a thread", func() {
			tid := ActiveBot.CreateNewThread("test-project", "go_developer", "valid-thread")
			Expect(tid).NotTo(BeEmpty())
		})

		It("can suggest a thread name based on the prompt", func() {
			prompt := "How do I implement a binary search tree in Go?"
			suggestedName, err := ActiveBot.SuggestThreadName(prompt)
			Expect(err).NotTo(HaveOccurred())
			Expect(suggestedName).NotTo(BeEmpty())

			// Check that the suggested name is not too long (should be 5 words or less)
			words := 0
			for _, c := range suggestedName {
				if c == ' ' {
					words++
				}
			}
			// Add 1 for the last word (which doesn't end with a space)
			words++
			Expect(words).To(BeNumerically("<=", 10), "Suggested name should be 10 words or less")
		})

		It("should automatically suggest a thread name when neither thread_id nor thread_name is provided", func() {
			prompt := "What is the best way to handle errors in Go?"
			request := completions.PromptRequest{
				Assistant:  "go_developer",
				ThreadId:   "",
				ThreadName: "", // Intentionally empty to trigger name suggestion
				Prompt:     prompt,
			}

			Eventually(func(g Gomega) {
				response, err := ActiveBot.QueryBot(&request)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(response).NotTo(BeEmpty())

				// The thread name should have been set automatically
				g.Expect(request.ThreadName).NotTo(BeEmpty())
				g.Expect(request.ThreadId).NotTo(BeEmpty())
			}, "2s", "2s").Should(Succeed())
		})
	})
})
