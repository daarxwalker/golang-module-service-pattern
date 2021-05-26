package core

import (
	"example/core/dd"
	"example/core/helper/tokenHelper"
)

func protect(provideService ProvideService, protectType string) bool {
	session, err := provideService.GetSession()
	if err != nil {
		dd.Print(err)
		return false
	}

	if provideService.GetRequest().Header.Get("X-Forwarded-For") != session.IP ||
		provideService.GetRequest().UserAgent() != session.Device {
		dd.Print("ip nebo device")
		return false
	}

	sessionToken, err := provideService.GetSessionToken()
	if err != nil {
		dd.Print(err)
		return false
	}

	token, err := tokenHelper.ValidateJWT(sessionToken)
	if err != nil {
		dd.Print(err)
		return false
	}

	if !token.Valid || (protectType == adminType && !session.Admin) {
		dd.Print("admin")
		return false
	}

	provideService.RenewSession()

	return true
}

func protectAny(provideService ProvideService) bool {
	session, err := provideService.GetSession()
	if err != nil {
		return false
	}

	if provideService.GetRequest().Header.Get("X-Forwarded-For") != session.IP ||
		provideService.GetRequest().UserAgent() != session.Device {
		return false
	}

	sessionToken, err := provideService.GetSessionToken()
	if err != nil {
		return false
	}

	token, err := tokenHelper.ValidateJWT(sessionToken)
	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	provideService.RenewSession()

	return true
}
