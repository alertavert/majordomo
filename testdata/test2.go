package testdata

import (
	"io"
	"os"
)

func anoter(f *os.File) ([]byte, error) {
	return io.ReadAll(f)
}
