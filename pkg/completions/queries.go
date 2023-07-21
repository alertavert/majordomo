/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package completions

import (
	"context"
	"fmt"
	"github.com/alertavert/gpt4-go/pkg/config"
	openai "github.com/sashabaranov/go-openai"
)

const (
	GoDeveloper = "go_developer"
	WebDesigner = "web_developer"

	MaxLogLen = 120
)

type PromptRequest struct {
	Prompt   string `json:"prompt"`
	Scenario string `json:"scenario"`
	Session  string `json:"session"`
}

// FIXME: Keeping the messages in memory is not a good idea.
var (
	userPrompts  []string
	botResponses []string

	// oaiClient is a singleton instance of the OpenAI client.
	oaiClient *openai.Client
)

// GetClient returns the singleton instance of the OpenAI client.
func GetClient() (*openai.Client, error) {
	if oaiClient == nil {
		cfg, err := config.LoadConfig()
		if err != nil {
			return nil, fmt.Errorf("error loading config: %w", err)
		}
		oaiClient = openai.NewClient(cfg.OpenAIApiKey)
	}
	return oaiClient, nil
}

func buildMessages(scenario string) ([]openai.ChatCompletionMessage, error) {
	messages := make([]openai.ChatCompletionMessage, 0, len(userPrompts)+len(botResponses)+1)
	if scenario == "" {
		return nil, fmt.Errorf("the scenario cannot be empty")
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: scenario,
	})

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
					Content: userPrompts[i],
				})
		}
	}
	return messages, nil
}

func QueryBot(prompt *PromptRequest) (string, error) {
	userPrompts = append(userPrompts, prompt.Prompt)
	scenario := GetScenarios().Scenarios[prompt.Scenario]
	messages, err := buildMessages(scenario)
	if err != nil {
		return "", err
	}
	client, err := GetClient()
	if err != nil {
		return "", err
	}
	fmt.Printf("Sending %d conversational items\n", len(messages))
	resp, err := client.CreateChatCompletion(
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
