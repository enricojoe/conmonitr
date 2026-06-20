package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// JWTAuth validates the Bearer access token from the Authorization header.
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := tokenFromRequest(r)
		if raw == "" {
			jsonError(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if _, err := parseToken(raw, "access"); err != nil {
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
	return ""
}

// wsAuthenticate reads the first WebSocket message as an access JWT and validates it.
// The client must send the token immediately after the connection is established.
// Returns false and closes the connection if auth fails.
func wsAuthenticate(conn *websocket.Conn) bool {
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, msg, err := conn.ReadMessage()
	conn.SetReadDeadline(time.Time{})
	if err != nil {
		conn.Close()
		return false
	}
	if _, err := parseToken(string(msg), "access"); err != nil {
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "unauthorized"))
		conn.Close()
		return false
	}
	return true
}
