/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package completions

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
)

const (
	GoDeveloper = "go_developer"
	WebDesigner = "web_developer"

	MaxLogLen = 120
)

type PromptRequest struct {
	// The user prompt.
	Prompt   string `json:"prompt"`
	// The scenario to use (selected by the user).
	Scenario string `json:"scenario"`
	// The session ID (if any) to keep track of past prompts/responses in the conversation.
	Session  string `json:"session,omitempty"`
	// The LLM model to use (selected by the user).
	Model    string `json:"model,omitempty"`
}

// FIXME: Keeping the messages in memory is not a good idea.
var (
	userPrompts  []string
	botResponses []string

	// oaiClient is a singleton instance of the OpenAI client.
	oaiClient *openai.Client
)

// SetClient configures the singleton instance of the OpenAI client.
func SetClient(client *openai.Client) {
	oaiClient = client
}

func BuildMessages(prompt *PromptRequest) ([]openai.ChatCompletionMessage, error) {
	messages := make([]openai.ChatCompletionMessage, 0,
		len(userPrompts) + len(botResponses) + 3)
	s := GetScenarios()
	if s == nil {
		return nil, fmt.Errorf("no scenarios found")
	}
	scenario, found := s.Scenarios[prompt.Scenario]
	if !found {
		return nil, fmt.Errorf("no scenario found for %s", prompt.Scenario)
	}
	// Common instructions for all scenarios.
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: s.GetCommon(),
	})
	// Scenario-specific instructions.
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: scenario,
	})
	// The user's prompt.
	userPrompts = append(userPrompts, prompt.Prompt)

	// FIXME: we should retrieve the stored conversation from the database.
	// The stored conversation thus far.
	for i := 0; i < len(userPrompts); i++ {
		messages = append(messages,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompts[i],
			})
		if i < len(botResponses) {
			messages = append(messages,
				openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: botResponses[i],
				})
		}
	}
	return messages, nil
}

func QueryBot(prompt *PromptRequest) (string, error) {
	messages, err := BuildMessages(prompt)
	if err != nil {
		return "", err
	}
	if oaiClient == nil {
		return "", fmt.Errorf("OpenAI client not initialized")
	}
	fmt.Printf("Sending %d conversational items\n", len(messages))
	resp, err := oaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			// FIXME: this should be configurable.
			Model:    openai.GPT4,
			Messages: messages,
		})

	if err != nil {
		return "", fmt.Errorf("error querying chatbot: %v", err)
	}

	if stopReason := resp.Choices[0].FinishReason; stopReason != "stop" {
		return "", fmt.Errorf("stopped for reason other than done: %s", stopReason)
	}
	botSays := resp.Choices[0].Message.Content
	botResponses = append(botResponses, botSays)

	totalTokens := resp.Usage.TotalTokens
	fmt.Printf("Tokens used: %d\n", totalTokens)
	logLen := len(botSays)
	if logLen > MaxLogLen {
		logLen = MaxLogLen
	}
	fmt.Printf("Bot says: %s\n", botSays[:logLen])
	return botSays, nil
}
