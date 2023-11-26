/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package completions

import (
	"context"
	"fmt"
	"github.com/alertavert/gpt4-go/pkg/config"
	"github.com/alertavert/gpt4-go/pkg/preprocessors"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"mime/multipart"
	"strings"
)

const (
	DefaultModel = openai.GPT4TurboPreview
)

type PromptRequest struct {
	// The user prompt.
	Prompt string `json:"prompt"`
	// The scenario to use (selected by the user).
	Scenario string `json:"scenario"`
	// The session ID (if any) to keep track of past prompts/responses in the conversation.
	Session string `json:"session,omitempty"`
}

type Majordomo struct {
	// The OpenAI Client
	Client *openai.Client
	// The Code Snippets CodeStore
	CodeStore preprocessors.CodeStoreHandler
	// The Model to use
	Model string
}

// NewMajordomo creates a new Majordomo instance from a Config struct.
func NewMajordomo(cfg *config.Config) (*Majordomo, error) {
	var assistant = new(Majordomo)
	assistant.Client = openai.NewClient(cfg.OpenAIApiKey)
	if assistant.Client == nil {
		return nil, fmt.Errorf("error initializing OpenAI client")
	}

	// Based on the active project, we set the code snippets directory.
	for _, p := range cfg.Projects {
		if p.Name == cfg.ActiveProject {
			destDir := strings.Join([]string{cfg.CodeSnippetsDir, p.Name}, "/")
			assistant.CodeStore = preprocessors.NewFilesystemStore(p.Location, destDir)
			log.Debug().
				Str("dest", destDir).
				Str("src", p.Location).
				Msg("code snippets filesystem store initialized")
			break
		}
		return nil, fmt.Errorf("no project found for %s", cfg.ActiveProject)
	}
	if cfg.Model == "" {
		assistant.Model = DefaultModel
	} else {
		assistant.Model = cfg.Model
	}
	log.Debug().
		Str("model", assistant.Model).
		Str("active_project", cfg.ActiveProject).
		Str("snippets", cfg.CodeSnippetsDir).
		Msg("assistant initialized")
	return assistant, nil
}

// FIXME: All the conversations should be indexed by the session ID, made up of the Project/Scenario
var (
	userPrompts  []string
	botResponses []string
)

// TODO: create a type to manage each conversation, and store them in a map.
// TODO: the key for the map should be the session ID, which is a combination of the project and scenario.
// TODO: the value should be a struct that holds the sequence of prompts and bot responses.

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

// FillPrompt fills the prompt with the code snippets.
func (m *Majordomo) FillPrompt(prompt *PromptRequest) error {
	p := prompt.Prompt
	oldLen := len(p)
	var parser = preprocessors.Parser{
		CodeMap: make(preprocessors.SourceCodeMap),
	}
	parser.ParsePrompt(p)
	err := m.CodeStore.GetSourceCode(&parser.CodeMap)
	if err != nil {
		log.Err(err).Msg("error retrieving source code")
		return err
	}
	prompt.Prompt, err = parser.FillPrompt(p)
	if err != nil {
		log.Err(err).Msg("error filling prompt")
		return err
	}
	log.Debug().
		Int("prompt_len", len(prompt.Prompt)).
		Int("old_len", oldLen).
		Int("code_snippets", len(parser.CodeMap)).
		Msg("filled prompt")
	return nil
}

// QueryBot queries the LLM with the given prompt.
func (m *Majordomo) QueryBot(prompt *PromptRequest) (string, error) {
	if m.Client == nil {
		return "", fmt.Errorf("OpenAI client not initialized")
	}
	if m.CodeStore == nil {
		return "", fmt.Errorf("code snippets store not initialized")
	}

	err := m.FillPrompt(prompt)
	if err != nil {
		return "", err
	}
	messages, err := BuildMessages(prompt)
	if err != nil {
		return "", err
	}
	log.Debug().
		Int("items", len(messages)).
		Str("scenario", prompt.Scenario).
		Str("session", prompt.Session).
		Str("model", m.Model).
		Msg("querying LLM")

	resp, err := m.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    m.Model,
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
	parser := preprocessors.Parser{
		CodeMap: make(preprocessors.SourceCodeMap),
	}
	err = parser.ParseBotResponse(botSays)
	if err != nil {
		return "", fmt.Errorf("error parsing bot response: %v", err)
	}
	err = m.CodeStore.PutSourceCode(parser.CodeMap)
	if err != nil {
		log.Err(err).Msg("error storing source code")
	}
	botResponses = append(botResponses, botSays)
	log.Debug().
		Int("tokens", resp.Usage.TotalTokens).
		Int("conversation_len", len(messages)).
		Int("response_len", len(botSays)).Send()
	return botSays, nil
}

func (m *Majordomo) SpeechToText(audioFile multipart.File) (string, error) {
	resp, err := m.Client.CreateTranscription(
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
