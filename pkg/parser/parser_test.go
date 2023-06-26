/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package parser_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alertavert/gpt4-go/pkg/parser"
)

var _ = Describe("ParseBotResponse", func() {
	It("should return an error when no valid code snippets are found", func() {
		err := parser.ParseBotResponse("This is not a code snippet")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no valid code snippets found"))
	})

	Context("with a well-formed code snippet", func() {
		BeforeEach(func() {
			err := parser.ParseBotResponse("'''test.txt\nThis is a test content\n'''")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err := os.Remove("test.txt")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should successfully write the correct content to the specified file path", func() {
			expectedContent := `This is a test content` + "\n"
			content, err := os.ReadFile("test.txt")
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal(expectedContent))
		})
	})

	Context("with malformed code snippets", func() {
		It("should return an error when missing closing triple-quotes", func() {
			err := parser.ParseBotResponse("'''test.txt\n missing closing triple-quotes")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no valid code snippets found"))
		})

		It("should return an error when missing opening triple-quotes", func() {
			err := parser.ParseBotResponse("missing opening triple-quotes\n'''")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no valid code snippets found"))
		})

		It("should return an error when the file path is malformed", func() {
			err := parser.ParseBotResponse("'''\nthis/is not: a /valid]path\nThis is content\n'''")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no valid code snippets found"))
		})
	})
})
