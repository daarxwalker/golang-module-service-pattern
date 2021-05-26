package passwordHelper

import (
	"github.com/matthewhartstonge/argon2"
)

func MakeArgonHash(value string) (string, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(value))
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func VerifyArgonHash(value string, hash string) bool {
	ok, err := argon2.VerifyEncoded([]byte(value), []byte(hash))
	return ok && err == nil
}
