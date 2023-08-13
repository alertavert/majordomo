/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package server

import (
	"github.com/alertavert/gpt4-go/pkg/preprocessors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

const DefaultFileLocation = "staging"

type Server struct {
	router *gin.Engine
	filestore preprocessors.CodeStoreHandler
}

var server *Server

func Setup(debug bool) *Server {
	// TODO: convert the server configuration to a struct
	srcDir, _ := os.Getwd()
	server = &Server{
		router: gin.Default(),
		// TODO: make the file locations configurable
		filestore: preprocessors.NewFilesystemStore(srcDir, DefaultFileLocation),
	}
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	server.router.Use(cors.Default())
	setupHandlers(server.router)
	return server
}

func (server *Server) Run(addr string) error {
	return http.ListenAndServe(addr, server.router)
}
