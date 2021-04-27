package util

import (
	"crypto/rand"
	"encoding/hex"
)

const idLength = 5

func GenerateID() (string, error) {
	rnd := make([]byte, idLength)
	_, err := rand.Read(rnd)

	return hex.EncodeToString(rnd), err
}

func GetId(URL string) string {
	return URL[len(URL)-idLength*2:]
}
