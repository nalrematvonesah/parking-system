package middleware

import (
	"context"
	"encoding/json"
	"gateway-service/internal/auth"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type ctxKey int

const (
	ctxUserID ctxKey = iota
)

func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(allowedOrigins))
	wildcard := false
	for _, o := range allowedOrigins {
		if o == "*" {
			wildcard = true
		}
		allowed[o] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if wildcard {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if origin != "" && allowed[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func JWTAuth(m *auth.Manager) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			header := r.Header.Get("Authorization")

			if !strings.HasPrefix(header, "Bearer ") {

				w.Header().Set("Content-Type", "application/json")

				w.WriteHeader(http.StatusUnauthorized)

				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "missing bearer token",
				})

				return
			}

			token := strings.TrimPrefix(header, "Bearer ")

			claims, err := m.Verify(token)

			if err != nil {

				w.Header().Set("Content-Type", "application/json")

				w.WriteHeader(http.StatusUnauthorized)

				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid token",
				})

				return
			}

			ctx := context.WithValue(
				r.Context(),
				ctxUserID,
				claims.UserID,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}

func UserIDFrom(ctx context.Context) int64 {
	v, _ := ctx.Value(ctxUserID).(int64)
	return v
}

func Recovery(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic in http handler",
						zap.String("path", r.URL.Path),
						zap.Any("panic", rec),
					)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func AccessLog(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := &statusWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(sw, r)
			log.Info("http",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", sw.status),
				zap.Duration("latency", time.Since(start)),
			)
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
