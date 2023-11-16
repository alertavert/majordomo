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
	"github.com/sashabaranov/go-openai"
)

// Release version of the server. It is expected to be set during build.
var Release = "UNKNOWN"

func main() {
	var port int
	var debug bool
	var model string
	var configPath string

	flag.IntVar(&port, "port", 8080, "Define the port the server will listen on for incoming requests")
	flag.BoolVar(&debug, "debug", false, "Set Debug log levels")
	flag.StringVar(&model, "model", "gpt-4", "Choose the LLM model to use")
	flag.StringVar(&configPath, "config", "", "Path to the configuration file; "+
		"if not specified, and the env var "+config.
		LocationEnv+" is not defined, it will use the default location: "+config.DefaultConfigLocation)
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
		Str("scenarios", cfg.ScenariosLocation).
		Str("snippets", cfg.CodeSnippetsDir).
		Msg("Loaded config")

	if err = completions.ReadScenarios(cfg.ScenariosLocation); err != nil {
		log.Fatal().Err(err).Msg("Error reading scenarios")
	}

	client := openai.NewClient(cfg.OpenAIApiKey)
	completions.SetClient(client)
	log.Info().Msg("OpenAI client configured")

	svr := server.Setup(debug)
	if svr == nil {
		log.Fatal().Msg("Error setting up server")
		return
	}
	log.Info().Msgf("Server configured & running on port %d", port)
	log.Fatal().Err(svr.Run(fmt.Sprintf(":%d", port))).
		Msg("Majordomo server exited")
}
