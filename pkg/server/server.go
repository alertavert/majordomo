/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	router *gin.Engine
}

var server *Server

func Setup(debug bool) *Server {
	server = &Server{
		router: gin.Default(),
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
