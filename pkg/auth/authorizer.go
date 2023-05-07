package auth

import (
	"crypto/sha512"
	"crypto/subtle"
	"net/http"

	"github.com/fi3te/monthly-planning-backend/pkg/config"
)

func IsAuthorized(cfg *config.Config, headers map[string]string) bool {
	r := createRequestForHeaders(headers)
	username, password, ok := r.BasicAuth()
	if !ok {
		return false
	}
	return constantTimeCompare(username, cfg.Username) && constantTimeCompare(password, cfg.Password)
}

func createRequestForHeaders(headers map[string]string) *http.Request {
	r, _ := http.NewRequest("", "", nil)
	for key, value := range headers {
		r.Header.Add(key, value)
	}
	return r
}

func constantTimeCompare(s1 string, s2 string) bool {
	s1Hash := sha512.Sum512([]byte(s1))
	s2Hash := sha512.Sum512([]byte(s2))
	return subtle.ConstantTimeCompare(s1Hash[:], s2Hash[:]) == 1
}
