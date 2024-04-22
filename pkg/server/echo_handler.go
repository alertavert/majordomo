/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package server

import (
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// echoPromptHandler is a simple handler that echoes back the prompt it receives
// after parsing it and substituting any code snippets.
func echoPromptHandler(m *completions.Majordomo) func(c *gin.Context) {
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
		err = m.FillPrompt(&requestBody)
		if err != nil {
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
