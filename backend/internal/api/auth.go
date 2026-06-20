package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   req.Username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	signed, err := token.SignedString(jwtSecret())
	if err != nil {
		jsonError(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{Token: signed})
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
