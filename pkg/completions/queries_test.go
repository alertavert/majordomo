/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package completions_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	openai "github.com/sashabaranov/go-openai"

	"github.com/alertavert/gpt4-go/pkg/completions"
)

const TestScenariosLocation = "../../testdata/test_scenarios.yaml"

type MockedClient struct{}

var noClientError = "OpenAI client not initialized"
var noScenarioError = "no scenario found for %s"

// TODO: we are not using Mocks for now, but we should.
func (m *MockedClient) CreateChatCompletion(ctx context.Context,
	req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	if req.Model != openai.GPT4 {
		return openai.ChatCompletionResponse{},
		fmt.Errorf("unexpected model: %s", req.Model)
	}
	choices := make([]openai.ChatCompletionChoice, 1)
	choices[0].Message.Content = "Hello"
	choices[0].FinishReason = "stop"
	chatCompletion := openai.ChatCompletionResponse{
		ID:      "123456",
		Object:  "test",
		Created: 0,
		Model:   req.Model,
		Choices: choices,
		Usage:   openai.Usage{TotalTokens: 123},
	}
	return chatCompletion, nil
}

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

	Context("QueryBot", func() {
		It("Should return 'OpenAI client not initialized' error", func() {
			_, err := completions.QueryBot(&req)
			Expect(err.Error()).To(Equal(noClientError))
		})
		// TODO: add more tests.
	})


	Context("BuildMessages function", func() {
		It("Should throw 'no scenario found' error", func() {
			req.Scenario = "unknown"
			_, err := completions.BuildMessages(&req)
			Expect(err.Error()).To(Equal(fmt.Sprintf(noScenarioError, "unknown")))
		})
		It("Should create messages when no error", func() {
			messages, err := completions.BuildMessages(&req)
			Expect(err).To(BeNil())
			Expect(len(messages)).To(Equal(3))
		})
		It("Should append Messages in the Correct Order", func() {
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
