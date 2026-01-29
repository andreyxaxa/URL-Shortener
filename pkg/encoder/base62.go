package encoder

import "strings"

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	var result strings.Builder

	base := int64(len(base62Chars))

	for num > 0 {
		result.WriteByte(base62Chars[num%base])
		num /= base
	}

	s := result.String()

	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func Decode(str string) int64 {
	var num int64
	base := int64(len(base62Chars))

	for _, c := range str {
		pos := strings.IndexRune(base62Chars, c)
		if pos == -1 {
			return 0
		}
		num = num*base + int64(pos)
	}

	return num
}
