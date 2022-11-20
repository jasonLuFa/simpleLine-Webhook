package util

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const alphabet_number = "abcdefghijklmnopqrstuvwxyz1234567890"

func RandomString(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomStringAndNumber(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		c := alphabet_number[rand.Intn(len(alphabet_number))]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
