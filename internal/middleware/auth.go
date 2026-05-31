package middleware

import (
	"errors"
	"net/http"
)

var AuthError = errors.New("Unauthorized")

func Authorize(r *http.Request, savedSessionToken string, savedCSRFToken string) error {
	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" || st.Value != savedSessionToken {
		return AuthError
	}

	csrf := r.Header.Get("X-CSRF-Token")
	if csrf == "" || csrf != savedCSRFToken {
		return AuthError
	}

	return nil
}
