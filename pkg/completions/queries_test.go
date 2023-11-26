/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package completions_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sashabaranov/go-openai"

	"github.com/alertavert/gpt4-go/pkg/completions"
)

const (
	TestScenariosLocation = "../../testdata/test_scenarios.yaml"
	NoClientError         = "OpenAI client not initialized"
	NoScenarioError       = "no scenario found for %s"
)

var _ = Describe("Building Prompts", func() {
	var (
		req completions.PromptRequest
	)
	BeforeEach(func() {
		Expect(completions.ReadScenarios(TestScenariosLocation)).ShouldNot(HaveOccurred())
		req = completions.PromptRequest{
			Prompt:   "prompt",
			Scenario: "test",
			Session:  "session",
		}
	})
	Context("BuildMessages function", func() {
		It("Should throw 'no scenario found' error", func() {
			req.Scenario = "unknown"
			_, err := completions.BuildMessages(&req)
			Expect(err.Error()).To(Equal(fmt.Sprintf(NoScenarioError, "unknown")))
		})
		It("Should create messages when no error", func() {
			completions.RemoveConversation()
			messages, err := completions.BuildMessages(&req)
			Expect(err).To(BeNil())
			Expect(len(messages)).To(Equal(3))
		})
		It("Should append Messages in the Correct Order", func() {
			completions.RemoveConversation()
			messages, err := completions.BuildMessages(&req)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(messages)).To(Equal(3))
			Expect(messages[0].Role).To(Equal(openai.ChatMessageRoleSystem))
			Expect(messages[1].Role).To(Equal(openai.ChatMessageRoleSystem))
			Expect(messages[2].Role).To(Equal(openai.ChatMessageRoleUser))
			Expect(messages[0].Content).To(Equal("common test scenario\n"))
			Expect(messages[1].Content).To(Equal("This is a test scenario\n"))
		})
		It("Should continue to append Messages in the Correct Order", func() {
			completions.RemoveConversation()
			for i := 0; i < 10; i++ {
				_, err := completions.BuildMessages(&req)
				Expect(err).ToNot(HaveOccurred())
			}
			req.Prompt = "Hello"
			messages, _ := completions.BuildMessages(&req)
			Expect(len(messages)).To(Equal(13))
			Expect(messages[12].Role).To(Equal(openai.ChatMessageRoleUser))
			Expect(messages[12].Content).To(Equal("Hello"))
		})
	})
})
