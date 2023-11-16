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
				_, err := config.LoadConfig("/tmp/config.yaml")
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with existing config path", func() {
			var location string
			BeforeEach(func() {
				location = "../../testdata/test_config.yaml"
			})
			It("should succeed", func() {
				c, err := config.LoadConfig(location)
				Expect(err).NotTo(HaveOccurred())
				Expect(c.OpenAIApiKey).To(Equal("test-key"))
				Expect(c.ScenariosLocation).To(HaveSuffix("test/scenarios.yaml"))
				Expect(c.CodeSnippetsDir).To(HaveSuffix("code/snippets"))
			})
			It("should expand relative paths", func() {
				c, err := config.LoadConfig(location)
				Expect(err).NotTo(HaveOccurred())
				Expect(c.ScenariosLocation).To(HavePrefix("../../testdata"))
				Expect(c.CodeSnippetsDir).To(HavePrefix("../../testdata"))
				Expect(c.CodeSnippetsDir).To(HaveSuffix("code/snippets"))
			})
		})
		Context("with configured projects", func() {
			var location string
			BeforeEach(func() {
				location = "../../testdata/test_config_projects.yaml"
			})
			It("should correctly load projects", func() {
				c, err := config.LoadConfig(location)
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Projects).To(HaveLen(2))

				project1 := c.Projects[0]
				Expect(project1.Name).To(Equal("test-project"))
				Expect(project1.Location).To(Equal("test/location"))
				Expect(project1.Description).To(Equal("test-description"))

				project2 := c.Projects[1]
				Expect(project2.Name).To(Equal("test-project-2"))
				Expect(project2.Location).To(Equal("test/location-2"))
				Expect(project2.Description).To(Equal("test-description-2"))
			})

		})
	})
	Describe("Save", func() {
		Context("with a valid configuration", func() {
			It("should successfully save the configuration as a yaml file", func() {
				c := &config.Config{
					OpenAIApiKey:      "test-key",
					ScenariosLocation: "test/scenarios.yaml",
					CodeSnippetsDir:   "code/snippets",
					Projects: []config.Project{
						{
							Name:        "test-project",
							Description: "test-description",
							Location:    "test/location",
						},
						{
							Name:        "test-project-2",
							Description: "test-description-2",
							Location:    "test/location-2",
						},
					},
				}

				filepath := os.TempDir() + "/config.yaml"
				err := c.Save(filepath)
				Expect(err).NotTo(HaveOccurred())

				content, err := os.ReadFile(filepath)
				Expect(err).NotTo(HaveOccurred())

				Expect(string(content)).To(ContainSubstring("test-key"))
				Expect(string(content)).To(ContainSubstring("test/scenarios.yaml"))
				Expect(string(content)).To(ContainSubstring("code/snippets"))
				Expect(string(content)).To(ContainSubstring("test-project"))
				Expect(string(content)).To(ContainSubstring("test-description"))
				Expect(string(content)).To(ContainSubstring("test/location"))
				Expect(string(content)).To(ContainSubstring("test-project-2"))
				Expect(string(content)).To(ContainSubstring("test-description-2"))
				Expect(string(content)).To(ContainSubstring("test/location-2"))
			})
		})

	})
})
