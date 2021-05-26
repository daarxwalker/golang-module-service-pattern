package tokenHelper

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var salt = os.Getenv("example_JWT_SALT")

func createRandomToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}

func CreateJWT(email string, admin bool) (string, error) {
	expires := 24 * time.Hour
	now := time.Now()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = email
	claims["createdAt"] = now
	claims["admin"] = admin
	claims["ExpiresAt"] = now.Add(expires).Unix()
	claims["NotBefore"] = now.Add(expires / 2).Unix()
	claims["IssuedAt"] = now.Unix()

	tokenString, err := token.SignedString([]byte(salt))

	return tokenString, err
}

func ValidateJWT(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("validate jwt failed")
		}
		return []byte(salt), nil
	})

	return token, err
}
