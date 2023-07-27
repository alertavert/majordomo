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

func promptHandler(c *gin.Context) {
	var requestBody completions.PromptRequest
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"response": "error",
			"message":  "No prompt received",
		})
		return
	}
	if requestBody.Scenario == "" {
		log.Debug().Msg("No scenario specified, using default")
		requestBody.Scenario = completions.GoDeveloper
	}

	botResponse, err := completions.QueryBot(&requestBody)
	if err != nil {
		log.Error().Err(err).Msg("Error querying bot")
		c.JSON(http.StatusInternalServerError, gin.H{
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
