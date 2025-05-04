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
		botResponse, err := m.QueryBot(&requestBody)
		if err != nil {
			log.Error().Err(err).Msg("Error querying bot")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		log.Debug().Msg("returning response")
		c.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"message":   botResponse,
			"thread_id": requestBody.ThreadId,
		})
	}
}
