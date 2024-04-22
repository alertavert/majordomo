/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package server

import (
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// global var to cache the Scenarios, to avoid reading the scenarios file on every request
//var cachedScenarios *completions.Scenarios

func scenariosHandler(c *gin.Context) {
	var scenarios = completions.GetScenarios()

	// return a response with all scenario titles
	c.JSON(http.StatusOK, gin.H{
		"status": "WARN -- this is deprecated and will be removed in a " +
			"future release. Use /assistants instead.",
		"scenarios": scenarios.GetScenarioNames(),
	})
}
