/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package server

import (
	"fmt"
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

var (
	ErrEmptyPrompt   = fmt.Errorf("prompt cannot be empty")
	ErrEmptyScenario = fmt.Errorf("scenario cannot be empty")
	ErrNoSession     = fmt.Errorf("session cannot be empty")
)

func ValidatePromptRequest(requestBody *completions.PromptRequest) error {
	if requestBody.Prompt == "" {
		return ErrEmptyPrompt
	}
	if requestBody.Scenario == "" {
		return ErrEmptyScenario
	}
	if requestBody.Session == "" {
		return ErrNoSession
	}
	return nil
}

func promptHandler(m *completions.Majordomo) func(c *gin.Context) {
	return func(c *gin.Context) {
		var requestBody completions.PromptRequest
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"response": "error",
				"message":  err.Error(),
			})
			return
		}
		if err := ValidatePromptRequest(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"response": "error",
				"message":  err.Error(),
			})
			return
		}
		botResponse, err := m.QueryBot(&requestBody)
		if err != nil {
			log.Error().Err(err).Msg("Error querying bot")
			c.JSON(http.StatusBadRequest, gin.H{
				"response": "error",
				"message":  err.Error(),
			})
			return
		}
		log.Debug().Msg("returning response")
		c.JSON(http.StatusOK, gin.H{
			"response": "success",
			"message":  botResponse,
		})
	}
}
