package ticket

import (
	cryptoRand "crypto/rand"
	"fmt"
	"io"
)

func (s service) generateSecureRandomString(length int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const maxIndex = len(letters)

	// Create a slice of bytes to hold the random indexes
	b := make([]byte, length)

	// Read random data from the crypto/rand source
	if _, err := io.ReadFull(cryptoRand.Reader, b); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}

	// Map the random bytes to the characters in the 'letters' string
	for i := 0; i < length; i++ {
		b[i] = letters[int(b[i])%maxIndex]
	}

	return string(b), nil
}
