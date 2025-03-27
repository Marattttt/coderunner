package auth

import (
	"crypto/rand"
	"math/big"
)

func GenerateAccessCode() (string) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 64

	result := make([]byte, length)
	for i := range result {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		result[i] = charset[index.Int64()]
	}

	return string(result)
}
