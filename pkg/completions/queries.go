/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package completions

import (
	"context"
	"errors"
	"fmt"
	"github.com/alertavert/gpt4-go/pkg/conversations"
	"github.com/go-playground/validator/v10"
	"mime/multipart"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"

	"github.com/alertavert/gpt4-go/pkg/config"
	"github.com/alertavert/gpt4-go/pkg/preprocessors"
)

const (
	DefaultModel = openai.GPT4Turbo
)

type PromptRequest struct {
	// The assistant to use (selected by the user); always required.
	Assistant string `json:"assistant" validate:"required"`

	// The Thread ID (if any) to keep track of past prompts/responses in the conversation.
	// If empty, a new conversation is started.
	ThreadId string `json:"thread_id,omitempty"`

	// The name of the thread, used to identify the conversation.
	// If both ThreadId and ThreadName are empty, a name will be suggested by the LLM.
	ThreadName string `json:"thread_name,omitempty"`

	// The user prompt.
	Prompt string `json:"prompt" validate:"required"`
}

// Validate checks if the PromptRequest has all required fields, using the validator package.
func (pr *PromptRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(pr); err != nil {
		// Convert validator errors to friendly messages
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors {
				return fmt.Errorf("%s field is required", e.Field())
			}
		}
		return err
	}
	return nil
}

type Majordomo struct {
	// The OpenAI Client
	Client *openai.Client

	// The Code Snippets CodeStore
	CodeStore preprocessors.CodeStoreHandler

	// Threads of conversation with the LLM model.
	Threads *conversations.ThreadStore

	// The Model to use
	Model string

	// The configuration object to manage the Projects in the server handlers
	Config *config.Config
}

// SuggestThreadName suggests a title for a thread based on the prompt text.
// It uses the OpenAI API to generate a title no longer than 5 words.
func (m *Majordomo) SuggestThreadName(prompt string) (string, error) {
	if m.Client == nil {
		return "", fmt.Errorf("OpenAI client not initialized")
	}

	ctx := context.Background()
	resp, err := m.Client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: m.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant that suggests concise titles for conversations. Provide a title that is no longer than 5 words based on the user's prompt. Return only the title, nothing else.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 20,
		},
	)
	if err != nil {
		log.Err(err).Msg("error suggesting thread name")
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no suggestions returned from OpenAI")
	}

	suggestedName := resp.Choices[0].Message.Content
	log.Debug().
		Str("suggested_name", suggestedName).
		Msg("suggested thread name")

	return suggestedName, nil
}

func NewMajordomo(cfg *config.Config) (*Majordomo, error) {
	var assistant = new(Majordomo)
	assistant.Client = openai.NewClient(cfg.OpenAIApiKey)
	if assistant.Client == nil {
		return nil, fmt.Errorf("error initializing OpenAI client")
	}

	// The LLM Model to use.
	if cfg.Model == "" {
		assistant.Model = DefaultModel
	} else {
		assistant.Model = cfg.Model
	}

	// Based on the active project, we set the code snippets directory.
	p := cfg.GetActiveProject()
	if p == nil {
		return nil, fmt.Errorf("no project found for %s", cfg.ActiveProject)
	}
	assistant.CodeStore = *preprocessors.GetCodeStoreHandler(p)
	assistant.Config = cfg
	assistant.Threads = conversations.NewThreadStore(cfg)
	if assistant.Threads == nil {
		return nil, fmt.Errorf("error initializing thread store")
	}

	log.Debug().
		Str("model", assistant.Model).
		Str("configured_snippets", cfg.CodeSnippetsDir).
		Str("threads_location", cfg.ThreadsLocation).
		Msg("Majordomo assistant initialized")
	log.Debug().
		Str("active_project", p.Name).
		Str("source_dir", p.Location).
		Str("code_snippets", p.ResolvedCodeSnippetsDir).
		Msg("Active Project set")
	return assistant, nil
}

func (m *Majordomo) SetActiveProject(projectName string) error {
	p := m.Config.GetProject(projectName)
	if p == nil {
		return fmt.Errorf("project %s not found", projectName)
	}
	m.Config.ActiveProject = projectName
	m.CodeStore = *preprocessors.GetCodeStoreHandler(p)
	log.Debug().
		Str("active_project", p.Name).
		Str("source_dir", p.Location).
		Str("code_snippets", p.ResolvedCodeSnippetsDir).
		Msg("New Active Project set")
	return nil
}

// PreparePrompt fills the prompt with the code snippets.
func (m *Majordomo) PreparePrompt(prompt *PromptRequest) error {
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

// CreateNewThread creates a new thread for the given project and returns the thread ID.
func (m *Majordomo) CreateNewThread(project, assistant, threadName string) string {
	t, err := m.Client.CreateThread(context.Background(), openai.ThreadRequest{
		Metadata: map[string]any{"project": project, "assistant": assistant, "thread_name": threadName},
	})
	if err != nil {
		log.Err(err).Msg("error creating thread")
		return ""
	}
	// TODO: we need to ask the LLM for a name and description.
	var newThread = conversations.Thread{
		ID:          t.ID,
		Name:        threadName,
		Assistant:   assistant,
		Description: "Some brief description for this thread",
	}
	err = m.Threads.AddThread(project, newThread)
	return t.ID
}

// QueryBot queries the LLM with the given prompt.
func (m *Majordomo) QueryBot(prompt *PromptRequest) (string, error) {
	if m.Client == nil {
		return "", fmt.Errorf("OpenAI client not initialized")
	}
	if m.CodeStore == nil {
		return "", fmt.Errorf("code snippets store not initialized")
	}

	err := m.PreparePrompt(prompt)
	if err != nil {
		return "", err
	}

	// TODO: create an appropriate context for the query.
	// Create a new conversation if the thread ID is empty.
	if prompt.ThreadId == "" {
		// If thread name is also empty, suggest a name based on the prompt
		if prompt.ThreadName == "" {
			suggestedName, err := m.SuggestThreadName(prompt.Prompt)
			if err != nil {
				log.Warn().
					Err(err).
					Msg("failed to suggest thread name, using default")
				prompt.ThreadName = "Untitled Conversation"
			} else {
				prompt.ThreadName = suggestedName
			}
			log.Debug().
				Str("thread_name", prompt.ThreadName).
				Msg("using suggested thread name")
		}

		log.Debug().
			Str("assistant", prompt.Assistant).
			Str("thread_name", prompt.ThreadName).
			Msg("creating new thread")
		prompt.ThreadId = m.CreateNewThread(m.Config.ActiveProject, prompt.Assistant, prompt.ThreadName)
	}
	log.Debug().
		Str("thread_id", prompt.ThreadId).
		Str("assistant", prompt.Assistant).
		Msg("thread ID set")
	// Creates a new conversation in the thread.
	msg, err := m.Client.CreateMessage(context.Background(), prompt.ThreadId,
		openai.MessageRequest{
			Role:    "user",
			Content: prompt.Prompt,
		})
	if err != nil {
		return "", err
	}
	log.Debug().
		// TODO: we should compute the number of tokens in debug mode only.
		Int("content_len", len(msg.Content)).
		Str("assistant", prompt.Assistant).
		Str("thread_id", prompt.ThreadId).
		Str("model", m.Model).
		Msg("querying LLM")

	// Find the assistant ID, given its name.
	if prompt.Assistant == "" {
		return "", fmt.Errorf("assistant name cannot be empty")
	}
	assistantId, err := m.GetAssistantId(prompt.Assistant)
	if err != nil {
		return "", fmt.Errorf("error getting assistant ID for '%s': %v", prompt.Assistant, err)
	}
	log.Debug().
		Str("assistant_id", assistantId).
		Str("assistant", prompt.Assistant).
		Msg("assistant found")
	// Create a Run - the model, and other parameters are set already in the Thread.
	run, err := m.Client.CreateRun(context.Background(), prompt.ThreadId, openai.RunRequest{
		// Model:       m.Model,
		AssistantID: assistantId,
	})
	if err != nil {
		return "", fmt.Errorf("error creating run: %v", err)
	}
	log.Debug().
		Str("run_id", run.ID).
		Str("thread_id", run.ThreadID).
		Str("assistant_id", run.AssistantID).
		Msg("created run")

	done := false
	// Get the response from the model.
	for !done {
		resp, err := m.Client.RetrieveRun(context.Background(), run.ThreadID, run.ID)
		if err != nil {
			return "", fmt.Errorf("error getting run: %v", err)
		}
		switch resp.Status {
		case openai.RunStatusInProgress, openai.RunStatusQueued:
			// TODO: we should have a configurable interval, or maybe exponential backoff.
			time.Sleep(5 * time.Second)
		case openai.RunStatusCompleted:
			log.Debug().
				Int("tokens", resp.Usage.TotalTokens).
				Msg("run completed")
			done = true
		case openai.RunStatusFailed:
			return "", fmt.Errorf("run failed: %v", resp.LastError.Message)
		case openai.RunStatusCancelled, openai.RunStatusCancelling, openai.RunStatusExpired:
			return "", fmt.Errorf("run cancelled or expired")
		case openai.RunStatusRequiresAction:
			log.Warn().
				Str("action", string(resp.RequiredAction.Type)).
				Msg("action required")
			return "", fmt.Errorf("action required")
		default:
			return "", fmt.Errorf("unexpected run status: %s", resp.Status)
		}
	}

	// Retrieve the most recent message in the Thread.
	messages, err := m.Client.ListMessage(
		context.Background(),
		prompt.ThreadId, nil, nil, nil, nil, nil)
	if err != nil {
		return "", fmt.Errorf("error listing messages: %v", err)
	}
	log.Debug().
		Int("messages", len(messages.Messages)).
		Str("last_id", *messages.LastID).
		Str("first_id", *messages.FirstID).
		Msg("messages")
	// TODO: should use the FirstID instead, and validate it's from `assistant`.
	botMessage := messages.Messages[0]
	// TODO: there is a lot more information in the response that we should log.
	if len(botMessage.Content) != 1 {
		log.Warn().
			Int("content_len", len(botMessage.Content)).
			Msg("unexpected content length")
	}
	botSays := botMessage.Content[0].Text.Value
	log.Debug().
		Str("bot_says", botSays).
		Msg("bot response")

	// Parse the response from the model.
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
	log.Debug().Msg("response parsed, code snippets stored")
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

// GetAssistantId returns the ID of the assistant with the given name.
// TODO: this should be cached somewhere, as the assistants change infrequently.
func (m *Majordomo) GetAssistantId(name string) (string, error) {
	ctx := context.Background()
	listAssistants, err := m.Client.ListAssistants(ctx, nil, nil, nil, nil)
	if err != nil {
		return "", fmt.Errorf("error listing assistants: %v", err)
	}
	for _, assistant := range listAssistants.Assistants {
		if assistant.Name == nil {
			log.Error().
				Str("assistant_id", assistant.ID).
				Msg("assistant has no name")
			continue
		}
		if *assistant.Name == name {
			return assistant.ID, nil
		}
	}
	return "", fmt.Errorf("assistant %s not found", name)
}

// CreateAssistants creates the OpenAI Assistants based on the instructions in the configuration file.
func (m *Majordomo) CreateAssistants(assistants *Assistants) error {
	ctx := context.Background()
	// TODO: This should be configurable.
	const DefaultTimeout = 1 * time.Second
	context.WithTimeout(ctx, DefaultTimeout)

	listAssistants, err := m.Client.ListAssistants(ctx, nil, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("error listing assistants: %v", err)
	}
	// Make it easier to check if an assistant already exists.
	existingAssistants := hashset.New()
	for _, assts := range listAssistants.Assistants {
		existingAssistants.Add(*assts.Name)
	}
	log.Debug().Msg(fmt.Sprintf("Existing Assistants: %v", existingAssistants.Values()))

	for name, instructions := range assistants.Instructions {
		if existingAssistants.Contains(name) {
			log.Warn().Str("assistant", name).Msg("assistant already exists, " +
				"updating not implemented yet")
			continue
		}
		inst := fmt.Sprintf("%s\n%s", assistants.Common, instructions)
		a, err := m.Client.CreateAssistant(ctx, openai.AssistantRequest{
			Model:        m.Model,
			Name:         &name,
			Description:  nil,
			Instructions: &inst,
		})
		if err != nil {
			log.Err(err).Str("assistant", name).Msg("error creating assistant")
			return err
		}
		log.Info().
			Str("assistant_id", a.ID).
			Str("assistant", *a.Name).
			Msg("assistant created")
	}
	return nil
}
