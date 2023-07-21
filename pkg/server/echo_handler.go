/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package server

import (
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)


func echoHandler(c *gin.Context) {
	var requestBody completions.PromptRequest
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"response": "error",
			"message":  err.Error(),
		})
		return
	}
	time.Sleep(1 * time.Second)
	c.JSON(http.StatusOK, gin.H{
		"response": "success",
		"message":  requestBody.Prompt,
	})
}
