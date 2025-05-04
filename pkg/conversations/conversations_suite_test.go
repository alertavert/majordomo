/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */
package conversations_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"testing"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	zerolog.SetGlobalLevel(zerolog.Disabled)
	RunSpecs(t, "Conversations Suite")
}
