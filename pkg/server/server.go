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
	// Health check
	r.GET("/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Prompt-related routes
	r.POST("/command", audioHandler(s.assistant))
	r.POST("/parse", parsePromptHandler(s.assistant))
	r.POST("/prompt", promptHandler(s.assistant))

	// Projects routes
	cfg := s.assistant.Config
	r.GET("/projects", projectsGetHandler(cfg))
	r.GET("/projects/:project_name", projectDetailsGetHandler(cfg))
	r.GET("/projects/:project_name/sessions", getSessionsForProjectHandler(s.assistant))
	r.POST("/projects", projectPostHandler(cfg))
	r.PUT("/projects", updateActiveProject(s.assistant))
	r.PUT("/projects/:project_name", projectPutHandler(cfg))
	r.DELETE("/projects/:project_name", projectDeleteHandler(cfg))

	// Names management routes
	r.GET("/assistants", assistantsGetHandler(s.assistant))
}

// SetupTestRoutes is a helper function to set up the routes for testing.
// Do not use this function in production code.
func SetupTestRoutes(r *gin.Engine, assistant *completions.Majordomo) {
	server := &Server{router: r, assistant: assistant}
	server.setupHandlers()
}
