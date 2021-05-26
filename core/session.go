package core

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-redis/redis/v8"
)

type Session struct {
	Email  string `json:"email"`
	ID     int    `json:"id"`
	Admin  bool   `json:"admin"`
	IP     string `json:"ip,omitempty"`
	Device string `json:"device,omitempty"`
}

func getSession(ctx context.Context, r *http.Request, cache *redis.Client) (*Session, error) {
	var session *Session

	cookie, err := r.Cookie(TokenKey)
	if err != nil {
		return nil, err
	}

	if cookie != nil && len(cookie.Value) > 0 {
		cachedSession := cache.Get(ctx, cookie.Value).Val()

		if len(cachedSession) == 0 {
			return nil, errors.New("session doesn't exist")
		}

		if err := json.Unmarshal([]byte(cachedSession), &session); err != nil {
			return nil, err
		}

		return session, nil
	}

	return nil, errors.New("cookie doesn't exist")
}

func getSessionByToken(ctx context.Context, token string, cache *redis.Client) (*Session, error) {
	var session *Session
	cachedSession := cache.Get(ctx, token).Val()
	if len(cachedSession) == 0 {
		return nil, errors.New("session doesn't exist or cookie is empty")
	}
	if err := json.Unmarshal([]byte(cachedSession), &session); err != nil {
		return nil, err
	}
	return session, nil
}

func getSessionToken(r *http.Request) (string, error) {
	token, err := r.Cookie(TokenKey)
	if err != nil {
		return "", err
	}

	return token.Value, nil
}

func getSessionString(ctx context.Context, r *http.Request, cache *redis.Client) (string, error) {
	cookie, err := r.Cookie(TokenKey)
	if err != nil {
		return "", err
	}

	if cookie != nil && len(cookie.Value) > 0 {
		cachedSession := cache.Get(ctx, cookie.Value).Val()

		if len(cachedSession) == 0 {
			return "", errors.New("session doesn't exist")
		}

		return cachedSession, nil
	}

	return "", errors.New("cookie doesn't exist")
}
