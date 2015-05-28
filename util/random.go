package util

import (
  "crypto/rand"
  "io"
)

const RANDOM_LENGTH int = 32

func Random () []byte {
  k := make([]byte, RANDOM_LENGTH)
  if _, err := io.ReadFull(rand.Reader, k); err != nil {
    return nil
  }
  return k
}
