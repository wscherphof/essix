package util

func NewToken() string {
	return string(URLEncode(Random()))
}
