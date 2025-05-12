/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package server

import (
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func promptHandler(m *completions.Majordomo) func(c *gin.Context) {
	return func(c *gin.Context) {
		var requestBody completions.PromptRequest
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			log.Error().Err(err).Msg("Cannot parse request body")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		if err := requestBody.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		lm := log.Debug().
			Str("assistant_name", requestBody.Assistant)
		var hasThreadId bool = false
		if requestBody.ThreadId != "" {
			hasThreadId = true
			lm.Str("thread_id", requestBody.ThreadId)
		}
		lm.Msg("Sending prompt to LLM")
		botResponse, err := m.QueryBot(&requestBody)
		if err != nil {
			log.Error().Err(err).Msg("Error querying bot")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		if !hasThreadId && requestBody.ThreadId != "" {
			log.Debug().
				Str("thread_id", requestBody.ThreadId).
				Msg("New thread created")
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"message":   botResponse,
			"thread_id": requestBody.ThreadId,
		})
	}
}
