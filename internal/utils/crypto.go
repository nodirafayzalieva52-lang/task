package utils

import "crypto/rand"

func GenerateRand() string {
	return rand.Text()
}
