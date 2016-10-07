package util

import (
	"encoding/base64"
)

func URLEncode(value []byte) []byte {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(encoded, value)
	return encoded
}

func URLDecode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
	b, err := base64.URLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:b], nil
}

func URLEncodeString(value string) string {
	return string(URLEncode([]byte(value)))
}

func URLDecodeString(value string) (string, error) {
	decoded, err := URLDecode([]byte(value))
	return string(decoded), err
}
