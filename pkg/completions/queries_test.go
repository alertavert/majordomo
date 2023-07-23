/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package completions_test

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alertavert/gpt4-go/pkg/completions"
)

type MockedClient struct{}

var noClientError = "OpenAI client not initialized"
var clientOpenAiError = "error querying chatbot: some error"
var noScenarioError = "no scenario found for %s"

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

var _ = Describe("Queries", func() {
	var (
		req completions.PromptRequest
	)

	BeforeEach(func() {
		req = completions.PromptRequest{
			Prompt:   "prompt",
			Scenario: completions.WebDesigner,
			Session:  "session",
		}
	})

	Context("QueryBot", func() {
		It("Should return 'OpenAI client not initialized' error", func() {
			_, err := completions.QueryBot(&req)
			Expect(err.Error()).To(Equal(noClientError))
		})

		It("Should return 'error querying chatbot: some error' for openai.GPT3", func() {
			req.Scenario = completions.GoDeveloper

			// create and set the mock client
			mockClient := MockedClient{}
			completions.SetClient(&mockClient)

			_, err := completions.QueryBot(&req)
			Expect(err.Error()).To(Equal(clientOpenAiError))
		})

		It("Should return 'stop' reason when no error", func() {
			// create and set the mock client
			mockClient := MockedClient{}
			completions.SetClient(&mockClient)

			response, err := completions.QueryBot(&req)
			Expect(err).To(BeNil())
			Expect(response).To(Equal("Hello"))
		})
	})

	Context("SetClient", func() {
		It("Should set the client", func() {
			// create and set the mock client
			mockClient := MockedClient{}
			completions.SetClient(&mockClient)

			_, err := completions.QueryBot(&req)
			Expect(err.Error()).NotTo(Equal(noClientError))
		})

		It("Should set the client to nil", func() {
			completions.SetClient(nil)

			_, err := completions.QueryBot(&req)
			Expect(err.Error()).To(Equal(noClientError))
		})
	})

	Context("buildMessages function", func() {
		It("Should throw 'no scenario found' error", func() {
			_, err := completions.BuildMessages(&req)
			Expect(err.Error()).To(Equal(fmt.Sprintf(noScenarioError, completions.WebDesigner)))
		})

		It("Should create messages when no error", func() {
			req.Scenario = completions.GoDeveloper
			messages, err := completions.BuildMessages(&req)
			Expect(err).To(BeNil())
			Expect(len(messages)).To(Equal(4))
		})
	})
})
