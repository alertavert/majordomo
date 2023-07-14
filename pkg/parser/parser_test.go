/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alertavert/gpt4-go/pkg/parser"
)

var _ = Describe("ParseBotResponse", func() {
	It("should return an error when no valid code snippets are found", func() {
		sourceCode, err := parser.ParseBotResponse("This is not a code snippet")
		Expect(err).To(HaveOccurred())
		Expect(sourceCode).To(BeNil())
		Expect(err.Error()).To(ContainSubstring("no valid code snippets found"))
	})

	Context("with a well-formed code snippet", func() {

		var source parser.SourceCode
		var err error

		BeforeEach(func() {
			source, err = parser.ParseBotResponse("'''test.txt\nThis is a test content\n'''")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should successfully extract the correct content to the source code map", func() {
			expectedContent := "This is a test content\n"
			Expect(len(source)).To(Equal(1))
			Expect(source["test.txt"]).To(Equal(expectedContent))
		})
	})

	Context("with malformed code snippets", func() {
		It("should return an error when missing closing triple-quotes", func() {
			source, err := parser.ParseBotResponse("'''test.txt\n missing closing triple-quotes")
			Expect(err).To(HaveOccurred())
			Expect(source).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("no valid code snippets found"))
		})

		It("should return an error when missing opening triple-quotes", func() {
			source, err := parser.ParseBotResponse("missing opening triple-quotes\n'''")
			Expect(err).To(HaveOccurred())
			Expect(source).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("no valid code snippets found"))
		})

		It("should return an error when the file path is malformed", func() {
			source, err := parser.ParseBotResponse("'''\nthis/is not: a /valid]path\nThis is content\n'''")
			Expect(err).To(HaveOccurred())
			Expect(source).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("no valid code snippets found"))
		})
	})
})
var _ = Describe("Parser", func() {
	Context("InsertSourceCode", func() {
		It("Should return the right file content", func() {
			text := "'''path/to/file.go\n'''"
			sourceCode := make(map[string]string)
			sourceCode["path/to/file.go"] = "this is the content"
			expectedResult := "'''path/to/file.go\nthis is the content'''"

			result, err := parser.InsertSourceCode(text, sourceCode)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedResult))
		})

		It("Should return error if no content found in SourceCode", func() {
			text := "'''path/to/file.go\n'''"
			sourceCode := make(map[string]string)

			_, err := parser.InsertSourceCode(text, sourceCode)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("no source code found for path: path/to/file.go"))
		})
	})
})
