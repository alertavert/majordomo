/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package server

import (
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	addr      string
	router    *gin.Engine
	assistant *completions.Majordomo
}

var server *Server

func NewServer(addr string, assistant *completions.Majordomo) *Server {
	server = &Server{
		router:    gin.Default(),
		assistant: assistant,
		addr:      addr,
	}
	server.router.Use(cors.Default())
	server.setupHandlers()
	return server
}

func (s *Server) SetDebugMode() {
	gin.SetMode(gin.DebugMode)
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.addr, s.router)
}

func (s *Server) setupHandlers() {
	r := s.router
	r.POST("/command", audioHandler(s.assistant))
	r.POST("/echo", echoHandler(s.assistant))
	r.POST("/prompt", promptHandler(s.assistant))
	r.GET("/scenarios", scenariosHandler)

	// Static routes
	r.Static("/web", "build/ui")
	r.Static("/static/css", "build/ui/static/css")
	r.Static("/static/js", "build/ui/static/js")
	r.Static("/static/media", "build/ui/static/media")
}
