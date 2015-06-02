package util

import (
	"crypto/rand"
	"io"
)

const length = 32

func Random() []byte {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}
