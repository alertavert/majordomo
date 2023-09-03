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

	// Static routes
	r.Static("/web", "build/ui")
	r.Static("/static/css", "build/ui/static/css")
	r.Static("/static/js", "build/ui/static/js")
	r.Static("/static/media", "build/ui/static/media")
}
