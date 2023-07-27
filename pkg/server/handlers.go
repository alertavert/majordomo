/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package server

import (
	"github.com/gin-gonic/gin"
)

func setupHandlers(r *gin.Engine) {
	r.POST("/command", audioHandler)
	r.POST("/echo", echoHandler)
	r.POST("/prompt", promptHandler)
	r.GET("/scenarios", scenariosHandler)
	r.Static("/web", "build/ui")
}
