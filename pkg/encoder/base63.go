package encoder

import (
	"strings"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

type Base63Encoder struct {
	len int
}

func NewBase63Encoder(len int) *Base63Encoder {
	return &Base63Encoder{len: len}
}

func (e *Base63Encoder) Encode(n int64) string {
	var tokenBuilder strings.Builder
	for n > 0 {
		idx := n % int64(len(chars))
		tokenBuilder.WriteByte(chars[idx])
		n /= int64(len(chars))
	}

	for tokenBuilder.Len() < e.len {
		tokenBuilder.WriteByte(chars[0])
	}

	return reverse(tokenBuilder.String())
}

func (e *Base63Encoder) Decode(token string) int64 {
	var id int64

	for _, char := range token {
		idx := strings.Index(chars, string(char))
		id = id*int64(len(chars)) + int64(idx)
	}
	return id
}

func reverse(s string) string {
	rns := []rune(s)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}
