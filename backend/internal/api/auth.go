package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTTL      = 15 * time.Minute
	refreshTTL     = 7 * 24 * time.Hour
	refreshCookie  = "conmonitr_refresh"
)

type tokenClaims struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	wantUser := os.Getenv("CONMONITR_USER")
	wantHash := os.Getenv("CONMONITR_PASSWORD_HASH")
	if wantUser == "" || wantHash == "" {
		jsonError(w, "server credentials not configured", http.StatusInternalServerError)
		return
	}

	if req.Username != wantUser || bcrypt.CompareHashAndPassword([]byte(wantHash), []byte(req.Password)) != nil {
		jsonError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	access, err := makeToken("access", req.Username, accessTTL)
	if err != nil {
		jsonError(w, "token generation failed", http.StatusInternalServerError)
		return
	}
	refresh, err := makeToken("refresh", req.Username, refreshTTL)
	if err != nil {
		jsonError(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	setRefreshCookie(w, r, refresh)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{Token: access})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(refreshCookie)
	if err != nil {
		jsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := parseToken(cookie.Value, "refresh")
	if err != nil {
		jsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	access, err := makeToken("access", claims.Subject, accessTTL)
	if err != nil {
		jsonError(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{Token: access})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookie,
		Value:    "",
		Path:     "/api/auth",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
	w.WriteHeader(http.StatusNoContent)
}

// makeToken creates a signed JWT with the given type, subject, and TTL.
func makeToken(tokenType, subject string, ttl time.Duration) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	return t.SignedString(jwtSecret())
}

// parseToken validates a JWT and checks that its type claim matches expectedType.
func parseToken(raw, expectedType string) (*tokenClaims, error) {
	var claims tokenClaims
	_, err := jwt.ParseWithClaims(raw, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if claims.Type != expectedType {
		return nil, fmt.Errorf("unexpected token type: %s", claims.Type)
	}
	return &claims, nil
}

func setRefreshCookie(w http.ResponseWriter, r *http.Request, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookie,
		Value:    value,
		Path:     "/api/auth",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(refreshTTL.Seconds()),
	})
}

func jwtSecret() []byte {
	if s := os.Getenv("CONMONITR_JWT_SECRET"); s != "" {
		return []byte(s)
	}
	return []byte("conmonitr-dev-secret-change-me")
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
