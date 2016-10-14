package util

import (
	"encoding/base64"
)

// URLEncode encodes a bite slice with base64.
func URLEncode(value []byte) []byte {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(encoded, value)
	return encoded
}

// URLDecode decodes a base64-encoded byte slice.
func URLDecode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
	b, err := base64.URLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:b], nil
}

// URLEncodeString encodes a string with base64.
func URLEncodeString(value string) string {
	return string(URLEncode([]byte(value)))
}

// URLDecode decodes a base64-encoded string.
func URLDecodeString(value string) (string, error) {
	decoded, err := URLDecode([]byte(value))
	return string(decoded), err
}
