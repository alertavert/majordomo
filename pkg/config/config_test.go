/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

// Author: M. Massenzio (marco@alertavert.com), 8/21/23

package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"

	"github.com/alertavert/gpt4-go/pkg/config"
)

var _ = Describe("Config", func() {
	Describe("LoadConfig", func() {
		Context("with non-existing config path", func() {
			It("should fail", func() {
				_, err := config.LoadConfig()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with existing config path", func() {
			BeforeEach(func() {
				err := os.Setenv("MAJORDOMO_CONFIG", "../../testdata/test_config.yaml")
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				err := os.Unsetenv("MAJORDOMO_CONFIG")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should succeed", func() {
				c, err := config.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(c.OpenAIApiKey).To(Equal("test-key"))
				Expect(c.ScenariosLocation).To(HaveSuffix("test/scenarios.yaml"))
				Expect(c.CodeSnippetsDir).To(HaveSuffix("code/snippets"))
			})
			It("should expand relative paths", func() {
				c, err := config.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(c.ScenariosLocation).To(HavePrefix("../../testdata"))
				Expect(c.CodeSnippetsDir).To(HavePrefix("../../testdata"))
			})
		})
	})
})
