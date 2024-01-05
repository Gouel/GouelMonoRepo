package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateSecretKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func main() {
	secretKey, err := generateSecretKey(64) // 64 octets pour une clé forte
	if err != nil {
		panic(err)
	}

	fmt.Println("JWT Secret Key:", secretKey)
}
