package util

// NewToken retuns a random 32 byte string.
func NewToken() string {
	return string(URLEncode(Random()))
}
