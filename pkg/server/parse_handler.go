/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package server

import (
	"fmt"
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// parsePromptHandler is a simple handler that echoes back the prompt it receives
// after parsing it and substituting any code snippets.
func parsePromptHandler(m *completions.Majordomo) func(c *gin.Context) {
	return func(c *gin.Context) {
		var requestBody completions.PromptRequest
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"response": "error",
				"message":  "Invalid JSON format: " + err.Error(),
			})
			return
		}

		// Validate required fields
		if err := requestBody.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"response": "error",
				"message":  fmt.Sprintf("Request has missing fields: %s", err.Error()),
			})
			return
		}
		if err := m.PreparePrompt(&requestBody); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"response": "error",
				"message":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": "success",
			"message":  requestBody.Prompt,
		})
	}
}
