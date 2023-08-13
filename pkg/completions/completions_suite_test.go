/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package completions_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCompletions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Completions Suite")
}
