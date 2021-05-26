package enviromentHelper

import (
	"os"
)

func IsDevelopment() bool {
	return os.Getenv("APP_ENV") == "development"
}

func IsTest() bool {
	return os.Getenv("TEST") == "true"
}

func GetEnviroment() string {
	if IsDevelopment() {
		return "development"
	}
	return "production"
}
