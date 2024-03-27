package activation

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateSecret(length int) string {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[seed.Intn(len(charset))]
	}

	return string(bytes)
}
