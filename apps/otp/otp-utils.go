package otp

import (
	"crypto/rand"
	"io"
)

func generateOTP(length int) string {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, length)

	_, err := io.ReadAtLeast(rand.Reader, b, length)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}

	return string(b)
}
