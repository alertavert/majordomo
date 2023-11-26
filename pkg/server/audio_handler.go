/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package server

import (
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func audioHandler(m *completions.Majordomo) func(c *gin.Context) {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("audio")
		if err != nil {
			log.Err(err).Msg("error getting audio content POSTed to /command")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// The name of the file is accessible as header.Filename
		log.Debug().
			Str("filename", header.Filename).
			Str("Content-Type", header.Header.Get("Content-Type")).
			Int("size", int(header.Size)).Msg("received audio file")

		text, err := m.SpeechToText(file)
		if err != nil {
			log.Err(err).Msg("error converting audio to text")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Debug().
			Str("text", text).
			Msg("converted audio to text")
		c.JSON(http.StatusOK, gin.H{
			"response": "success",
			"message":  text,
		})
	}
}
