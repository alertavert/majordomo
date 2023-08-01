/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package completions

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
	"mime/multipart"
)

const (
	GoDeveloper = "go_developer"
	WebDesigner = "web_developer"

	MaxLogLen = 120

	DefaultModel = openai.GPT4
)

type PromptRequest struct {
	// The user prompt.
	Prompt string `json:"prompt"`
	// The scenario to use (selected by the user).
	Scenario string `json:"scenario"`
	// The session ID (if any) to keep track of past prompts/responses in the conversation.
	Session string `json:"session,omitempty"`
	// The LLM model to use (selected by the user).
	Model string `json:"model,omitempty"`
}

// FIXME: Keeping the messages in memory is not a good idea.
var (
	userPrompts  []string
	botResponses []string

	// oaiClient is an instance of the OpenAI client.
	oaiClient *openai.Client
)

// SetClient configures the singleton instance of the OpenAI client.
func SetClient(client *openai.Client) {
	oaiClient = client
}

func BuildMessages(prompt *PromptRequest) ([]openai.ChatCompletionMessage, error) {
	messages := make([]openai.ChatCompletionMessage, 0,
		len(userPrompts)+len(botResponses)+3)
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
	if prompt.Model == "" {
		log.Debug().Msg("using default model")
		prompt.Model = DefaultModel
	}
	log.Debug().
		Int("items", len(messages)).
		Str("model", prompt.Model).
		Msg("querying LLM")

	resp, err := oaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    prompt.Model,
			Messages: messages,
		})
	if err != nil {
		return "", fmt.Errorf("error querying chatbot: %v", err)
	}
	stopReason := resp.Choices[0].FinishReason
	if stopReason != "stop" {
		if stopReason == "length" {
			log.Debug().Msg("too many tokens, response truncated")
			botResponses = botResponses[1:]
			userPrompts = userPrompts[1:]
			return "", fmt.Errorf("too many tokens (%d), "+
				"dropped older conversations. Please re-send the request, "+
				"or consider starting a new conversation", resp.Usage.TotalTokens)
		}
		return "", fmt.Errorf("stopped for reason other than done: %s", stopReason)
	}
	botSays := resp.Choices[0].Message.Content
	botResponses = append(botResponses, botSays)
	log.Debug().
		Int("tokens", resp.Usage.TotalTokens).
		Int("conversation_len", len(messages)).
		Int("response_len", len(botSays)).Send()
	return botSays, nil
}

func SpeechToText(audioFile multipart.File) (string, error) {
	resp, err := oaiClient.CreateTranscription(
		context.Background(),
		openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: "audio.mp3",
			Reader:   audioFile,
			Format:   openai.AudioResponseFormatText,
		})
	if err != nil {
		return "", fmt.Errorf("error converting audio to text: %v", err)
	}
	return resp.Text, nil
}

// RemoveConversation deleted all past responses and prompts
func RemoveConversation() {
	userPrompts = userPrompts[:0]
	botResponses = botResponses[:0]
}
