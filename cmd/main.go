/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alertavert/gpt4-go/pkg/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Release version of the server. It is expected to be set during build.
var Release = "UNKNOWN"

func main() {
	var port int
	var debug bool
	var model string

	// Read the command-line flags
	flag.IntVar(&port, "port", 8080, "Define the port the server will listen on for incoming requests")
	flag.BoolVar(&debug, "debug", false, "Set Debug log levels")
	flag.StringVar(&model, "model", "gpt-4", "Choose the LLM model to use")

	// Parse the command-line flags
	flag.Parse()

	// Set up logging
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg(fmt.Sprintf("Starting >>> Majordomo Server <<< Rel. %s >>>", Release))

	if err := server.Setup(); err != nil {
		log.Fatal().Err(err).Msg("Error setting up server")
		return
	}
	log.Fatal().Err(server.Run(fmt.Sprintf(":%d", port))).Msg("Error running server")
}
