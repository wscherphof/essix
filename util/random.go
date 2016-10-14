/*
Package util provides some utility functions.
*/
package util

import (
	"crypto/rand"
	"io"
)

const length = 32

// Random returns a 32 random bytes.
func Random() []byte {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}
