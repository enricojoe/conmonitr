package api

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth validates the Bearer token from Authorization header or ?token= query param.
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := tokenFromRequest(r)
		if raw == "" {
			jsonError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		_, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret(), nil
		})
		if err != nil {
			jsonError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func tokenFromRequest(r *http.Request) string {
	if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
		return h[7:]
	}
	return r.URL.Query().Get("token")
}
