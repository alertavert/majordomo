/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/alertavert/gpt4-go/pkg/config"
	"github.com/alertavert/gpt4-go/pkg/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Release version of the server. It is expected to be set during build.
var Release = "UNKNOWN"

func main() {
	var port int
	var debug bool
	var configPath string
	var shouldCreateAssistants bool

	flag.IntVar(&port, "port", 8080, "Define the port the server will listen on for incoming requests")
	flag.BoolVar(&debug, "debug", false, "Set Debug log levels")
	flag.StringVar(&configPath, "config", "", "Path to the configuration file; "+
		"if not specified, and the env var "+config.
		LocationEnv+" is not defined, it will use the default location: "+config.DefaultConfigLocation)
	flag.BoolVar(&shouldCreateAssistants, "create", false, "Create the OpenAI Assistants")
	flag.Parse()

	// Set up logging
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Info().Msg(fmt.Sprintf("Starting >>> Majordomo Server <<< Rel. %s >>>", Release))

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}
	log.Info().
		Str("configPath", cfg.LoadedFrom).
		Str("instructions", cfg.AssistantsLocation).
		Msg("Loaded config")


	majordomo, err := completions.NewMajordomo(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing Majordomo")
	}
	if shouldCreateAssistants {
		assistants, err := completions.ReadInstructions(cfg.AssistantsLocation)
		if err != nil {
			log.Fatal().Err(err).Msg("Error reading scenarios")
		}
		log.Info().Msg("Creating OpenAI Assistants")
		if err = majordomo.CreateAssistants(assistants); err != nil {
			log.Fatal().Err(err).Msg("Error creating OpenAI Assistants")
		}
		log.Info().Msg("OpenAI Assistants created")
	}
	log.Info().Msg("Majordomo initialized, starting server")
	svr := server.NewServer(fmt.Sprintf(":%d", port), majordomo)
	if svr == nil {
		log.Fatal().Msg("Error setting up server")
		return
	}
	if debug {
		svr.SetDebugMode()
	}
	log.Info().Msgf("Server configured & running on port %d", port)
	log.Fatal().Err(svr.Run()).
		Msg("Majordomo server exited")
}
