package kolkata

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func hashString(secret string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), 10)
	if err == nil {
		return hash
	}
	return []byte{} //(secret)
}

func randomCharacters(length_opt ...int) string {

	length := 16
	if len(length_opt) > 0 {
		length = length_opt[0]
	}
	rb := make([]byte, length)
	cryptorand.Read(rb)
	return base64.URLEncoding.EncodeToString(rb)
}
