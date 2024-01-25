package server_test

import (
	"io"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

// MkTempConfigFile creates a temporary copy of the test config file, and returns its path.
func MkTempConfigFile(src string) (dest string, err error) {
	var sourceFile *os.File
	sourceFile, err = os.Open(src)
	if err != nil {
		return
	}
	defer sourceFile.Close()

	dest = os.TempDir() + "/test_config.yaml"
	destFile, err := os.Create(dest)
	if err != nil {
		return
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return
	}
	return
}
