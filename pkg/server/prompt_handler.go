/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package server

import (
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func promptHandler(c *gin.Context) {
	var requestBody struct {
		Prompt string `json:"prompt"`
	}
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"response": "error",
			"message":  "No prompt received",
		})
		return
	}
	botResponse, err := completions.QueryBot(requestBody.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": "success",
		"message":  botResponse,
	})
}
