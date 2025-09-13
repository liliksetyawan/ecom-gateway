package middleware

import (
	"context"
	"ecom-gateway/server"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type ContextData struct {
	UserID string
}

type key int

const ContextKey key = 0

type Middleware struct {
	redis *server.RedisClient
}

func NewMiddleware(redis *server.RedisClient) *Middleware {
	return &Middleware{redis: redis}
}

func (m *Middleware) LoggingAndAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		authHeader := r.Header.Get("Authorization")
		var userID string
		var err error
		var ctxData *ContextData

		if !strings.Contains(r.URL.Path, "/register") && !strings.Contains(r.URL.Path, "/login") {
			//if authHeader != "" {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err = m.redis.ValidateToken(token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				logJSON(r, start, http.StatusUnauthorized, "unauthorized")
				return
			}

			ctxData = &ContextData{
				UserID: userID,
			}
			//}
		}

		ctx := context.WithValue(r.Context(), ContextKey, ctxData)
		next.ServeHTTP(w, r.WithContext(ctx))

		logJSON(r, start, http.StatusOK, userID)
	})
}

func logJSON(r *http.Request, start time.Time, status int, user string) {
	logData := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"method":    r.Method,
		"path":      r.URL.Path,
		"status":    status,
		"latency":   time.Since(start).String(),
		"user":      user,
	}
	logBytes, _ := json.Marshal(logData)
	log.Println(string(logBytes))
}
