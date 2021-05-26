package core

import (
	"fmt"
	"net/http"
	"os"

	"example/core/helper/enviromentHelper"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
)

/*
 *
 * CORS MIDDLEWARE
 *
 */
func makeCORS() *cors.Cors {
	if enviromentHelper.IsDevelopment() {
		return cors.New(cors.Options{
			AllowedOrigins:   []string{fmt.Sprintf("http://%s.loc", os.Getenv("DOMAIN_MAIN"))},
			AllowCredentials: true,
			AllowedMethods:   []string{"POST"},
			Debug:            true,
		})
	}

	return cors.New(cors.Options{
		AllowedOrigins:   []string{fmt.Sprintf("http://%s.%s", os.Getenv("DOMAIN_MAIN"), os.Getenv("DOMAIN_TLD"))},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST"},
		Debug:            false,
	})
}

func corsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return makeCORS().Handler(next)
	}
}

/*
 *
 * SECURE MIDDLEWARE
 *
 */
func makeSecure() *secure.Secure {
	if enviromentHelper.IsDevelopment() {
		return secure.New(secure.Options{
			IsDevelopment:           true,
			AllowedHosts:            []string{fmt.Sprintf("%s:%s", os.Getenv("example_NAME"), os.Getenv("example_PORT")), fmt.Sprintf("%s.loc", os.Getenv("DOMAIN_MAIN"))},
			SSLProxyHeaders:         map[string]string{"X-Forwarded-Proto": "http"},
			CustomFrameOptionsValue: "SAMEORIGIN",
			ContentTypeNosniff:      true,
			BrowserXssFilter:        true,
			ReferrerPolicy:          "same-origin",
			FeaturePolicy:           "vibrate 'none';",
		})
	}
	return secure.New(secure.Options{
		IsDevelopment:           false,
		AllowedHosts:            []string{fmt.Sprintf("%s:%s", os.Getenv("example_NAME"), os.Getenv("example_PORT")), fmt.Sprintf("%s.%s", os.Getenv("DOMAIN_MAIN"), os.Getenv("DOMAIN_TLD"))},
		HostsProxyHeaders:       []string{"X-Forwarded-Hosts"},
		SSLRedirect:             true,
		SSLHost:                 "",
		SSLProxyHeaders:         map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:              31536000,
		STSIncludeSubdomains:    true,
		STSPreload:              true,
		ForceSTSHeader:          true,
		FrameDeny:               true,
		CustomFrameOptionsValue: "SAMEORIGIN",
		ContentTypeNosniff:      true,
		BrowserXssFilter:        true,
		ContentSecurityPolicy:   "default-src 'self'",
		ReferrerPolicy:          "same-origin",
		FeaturePolicy:           "vibrate 'none';",
	})
}

func secureMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return makeSecure().Handler(next)
	}
}

/*
 *
 * RATE LIMIT MIDDLEWARE
 *
 */
func rateLimitMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return getRateLimiter(next)
	}
}
