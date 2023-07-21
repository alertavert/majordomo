/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package server

import (
	"fmt"
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"net/http"

	"github.com/alertavert/gpt4-go/pkg/config"
)

type Server struct {
	oaiClient *openai.Client
	router *gin.Engine
	config *config.Config
}

var server *Server

func Setup() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	client := openai.NewClient(cfg.OpenAIApiKey)
	server = &Server{
		oaiClient: client,
		router: gin.Default(),
		config: &cfg,
	}
	server.router.Use(cors.Default())
	setupHandlers(server.router)
	err = completions.ReadScenarios(cfg.ScenariosLocation)
	return err
}

func Run(addr string) error {
	return http.ListenAndServe(addr, server.router)
}
