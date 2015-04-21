package util

import (
  "encoding/base64"
  "crypto/rand"
  "io"
)

func URLEncode (value []byte) []byte {
  encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
  base64.URLEncoding.Encode(encoded, value)
  return encoded
}

func URLDecode (value []byte) ([]byte, error) {
  decoded := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
  b, err := base64.URLEncoding.Decode(decoded, value)
  if err != nil {
    return nil, err
  }
  return decoded[:b], nil
}

func Random (length int) []byte {
  k := make([]byte, length)
  if _, err := io.ReadFull(rand.Reader, k); err != nil {
    return nil
  }
  return k
}
