package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
)

// charset is the set of characters that are used to generate random passwords.
//
//nolint:gochecknoglobals  // This is a constant.
var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(length int) (string, error) {
	buffer := make([]rune, length)

	//nolint:varnamelen // i is a common variable name for loops.
	for i := range buffer {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("rand.Int err: %w", err)
		}

		buffer[i] = charset[n.Int64()]
	}

	return string(buffer), nil
}

// String generates a random string of the given length.
// NOTE: This function is not safe for use with passwords.
func String(length int) string {
	randomBytes := make([]rune, length)

	// Fill the byte slice with random characters from the charset
	for i := range randomBytes {
		//nolint:gosec  // This is not a security-sensitive operation.
		randomBytes[i] = charset[mrand.Intn(len(charset))]
	}

	// Convert the byte slice to a string and return
	return string(randomBytes)
}
