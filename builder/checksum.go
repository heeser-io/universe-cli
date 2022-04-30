package builder

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func Checksum(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
