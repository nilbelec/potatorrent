package util

func ToUtf8(text string) string {
	buffer := []byte(text)
	buf := make([]rune, len(buffer))
	for i, b := range buffer {
		buf[i] = rune(b)
	}
	return string(buf)
}
