/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package completions

import (
	"fmt"
	"github.com/sashabaranov/go-openai"
)

type Session struct {
	SessionID  string
	ScenarioID string

	initialized   bool
	systemPrompts []openai.ChatCompletionMessage
	userPrompts   []openai.ChatCompletionMessage
	botResponses  []openai.ChatCompletionMessage
}

// AddPrompt adds a user prompt to the session.
func (s *Session) AddPrompt(prompt string) {
	s.userPrompts = append(s.userPrompts, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})
}

// AddResponse adds a bot response to the session.
func (s *Session) AddResponse(response string) {
	s.botResponses = append(s.botResponses, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: response,
	})
}

// IsEmpty returns true if the session has no prompts or responses.
func (s *Session) IsEmpty() bool {
	return len(s.userPrompts) == 0 && len(s.botResponses) == 0
}

// Init initializes a new session with the given scenario, adding both common
// and scenario-specific instructions.
func (s *Session) Init(scenarioId string) error {
	if s.initialized {
		return fmt.Errorf("session %s already has an ongoing conversation", s.SessionID)
	}
	sc := GetScenarios()
	if sc == nil {
		return fmt.Errorf("no scenarios found")
	}
	scenario, found := sc.Scenarios[scenarioId]
	if !found {
		return fmt.Errorf("no scenario %s found", scenarioId)
	}
	s.ScenarioID = scenarioId
	// Common instructions for all scenarios.
	s.systemPrompts = append(s.systemPrompts, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: sc.GetCommon(),
	})
	// Scenario-specific instructions.
	s.systemPrompts = append(s.systemPrompts, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: scenario,
	})
	s.initialized = true
	return nil
}

// Clip removes the oldest n items from the session.
func (s *Session) Clip(n int) {
	if n <= 0 {
		return
	}
	if n > len(s.userPrompts) {
		n = len(s.userPrompts)
	}
	s.userPrompts = s.userPrompts[n:]
	s.botResponses = s.botResponses[n:]
}

// Clear removes all items from the session.
func (s *Session) Clear() {
	s.userPrompts = make([]openai.ChatCompletionMessage, 0)
	s.botResponses = make([]openai.ChatCompletionMessage, 0)
}

// GetConversation returns the conversation so far, by interleaving the user prompts
// and bot responses.
// The conversation is returned as a slice of ChatCompletionMessage structs, with the
// scenario instructions at the beginning.
func (s *Session) GetConversation() []openai.ChatCompletionMessage {
	var conversation []openai.ChatCompletionMessage
	conversation = append(conversation, s.systemPrompts...)
	for i, prompt := range s.userPrompts {
		conversation = append(conversation, prompt)
		if i < len(s.botResponses) {
			conversation = append(conversation, s.botResponses[i])
		}
	}
	return conversation
}

// GetUserPrompts returns the user prompts so far.
// Should be used for testing only.
func (s *Session) GetUserPrompts() []openai.ChatCompletionMessage {
	return s.userPrompts
}

// GetBotResponses returns the bot responses so far.
// Should be used for testing only.
func (s *Session) GetBotResponses() []openai.ChatCompletionMessage {
	return s.botResponses
}

func NewSession(sessionId string) *Session {
	if sessionId == "" {
		return nil
	}
	return &Session{
		SessionID:     sessionId,
		initialized:   false,
		systemPrompts: make([]openai.ChatCompletionMessage, 0),
		userPrompts:   make([]openai.ChatCompletionMessage, 0),
		botResponses:  make([]openai.ChatCompletionMessage, 0),
	}
}
