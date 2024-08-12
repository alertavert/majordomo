/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"

	"github.com/alertavert/gpt4-go/pkg/config"
)

const (
	testConfigLocation         = "../../testdata/test_config.yaml"
	testConfigProjectsLocation = "../../testdata/test_config_projects.yaml"
)

var _ = Describe("Config", func() {
	Describe("LoadConfig", func() {
		Context("with non-existing config path", func() {
			It("should fail", func() {
				_, err := config.LoadConfig("/etc/fake-config.yaml")
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with existing config path", func() {
			It("should succeed", func() {
				c, err := config.LoadConfig(testConfigLocation)
				Expect(err).NotTo(HaveOccurred())
				Expect(c.OpenAIApiKey).To(Equal("test-key"))
				Expect(c.AssistantsLocation).To(HaveSuffix("test/assistants.yaml"))
				Expect(c.CodeSnippetsDir).To(Equal(".majordomo"))
			})
			It("should expand relative paths", func() {
				c, err := config.LoadConfig(testConfigLocation)
				Expect(err).NotTo(HaveOccurred())
				Expect(c.AssistantsLocation).To(HavePrefix("../../testdata"))
				Expect(c.CodeSnippetsDir).To(Equal(".majordomo"))
			})
		})
		Context("with configured projects", func() {
			It("should correctly load projects", func() {
				c, err := config.LoadConfig(testConfigProjectsLocation)
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
					AssistantsLocation: "test/assistants.yaml",
					CodeSnippetsDir:   ".snippets",
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

				content, err := config.LoadConfig(filepath)
				Expect(err).NotTo(HaveOccurred())

				Expect(content.OpenAIApiKey).To(Equal("test-key"))
				Expect(content.AssistantsLocation).To(HaveSuffix("test/assistants.yaml"))
				Expect(content.CodeSnippetsDir).To(HaveSuffix("snippets"))
				Expect(content.Projects).To(HaveLen(2))
			})

		})
	})
})
