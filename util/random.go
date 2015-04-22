package util

import (
  "crypto/rand"
  "io"
)

func Random (length int) []byte {
  k := make([]byte, length)
  if _, err := io.ReadFull(rand.Reader, k); err != nil {
    return nil
  }
  return k
}
