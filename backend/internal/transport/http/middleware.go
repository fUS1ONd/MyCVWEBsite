package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"personal-web-platform/config"
	"personal-web-platform/internal/domain"

	"github.com/go-chi/chi/v5/middleware"
)

type contextKey string

const userContextKey contextKey = "user"

// AuthRequired middleware checks if user is authenticated
func (h *Handler) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := h.getSessionToken(r)
		if token == "" {
			RespondUnauthorized(w, "authentication required")
			return
		}

		user, err := h.services.Auth.ValidateSession(r.Context(), token)
		if err != nil {
			h.log.Error("failed to validate session", "error", err)
			RespondInternalError(w)
			return
		}

		if user == nil {
			RespondUnauthorized(w, "invalid or expired session")
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminRequired middleware checks if user is admin
func (h *Handler) AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := h.getUserFromContext(r.Context())
		if user == nil {
			RespondUnauthorized(w, "authentication required")
			return
		}

		if user.Role != domain.RoleAdmin {
			RespondForbidden(w, "admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getUserFromContext extracts user from request context
func (h *Handler) getUserFromContext(ctx context.Context) *domain.User {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	if !ok {
		return nil
	}
	return user
}

// RequestLogger is a custom logging middleware that logs HTTP requests
func (h *Handler) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			duration := time.Since(start)
			h.log.Info("http request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.Duration("duration", duration),
				slog.String("user_agent", r.UserAgent()),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}

// CORS is a middleware that handles CORS headers
func (h *Handler) CORS(cfg *config.CORS) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !cfg.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range cfg.AllowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}

			if allowed {
				if origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				} else if len(cfg.AllowedOrigins) > 0 && cfg.AllowedOrigins[0] == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				}

				if cfg.AllowedCreds {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ","))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ","))

				if len(cfg.ExposedHeaders) > 0 {
					w.Header().Set("Access-Control-Expose-Headers", strings.Join(cfg.ExposedHeaders, ","))
				}

				if cfg.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", cfg.MaxAge))
				}
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// rateLimiter holds rate limiting state
type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// visitor tracks requests from a single IP
type visitor struct {
	lastSeen time.Time
	count    int
	resetAt  time.Time
}

// newRateLimiter creates a new rate limiter
func newRateLimiter(requestsLimit int, windowSeconds int) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		limit:    requestsLimit,
		window:   time.Duration(windowSeconds) * time.Second,
	}

	// Cleanup old visitors every minute
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			rl.cleanupOldVisitors()
		}
	}()

	return rl
}

// cleanupOldVisitors removes visitors that haven't been seen in a while
func (rl *rateLimiter) cleanupOldVisitors() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, v := range rl.visitors {
		if now.Sub(v.lastSeen) > rl.window*2 {
			delete(rl.visitors, ip)
		}
	}
}

// allow checks if a request from the given IP is allowed
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, exists := rl.visitors[ip]

	if !exists {
		rl.visitors[ip] = &visitor{
			lastSeen: now,
			count:    1,
			resetAt:  now.Add(rl.window),
		}
		return true
	}

	v.lastSeen = now

	// Reset counter if window has passed
	if now.After(v.resetAt) {
		v.count = 1
		v.resetAt = now.Add(rl.window)
		return true
	}

	// Increment counter
	v.count++

	return v.count <= rl.limit
}

// RateLimit is a middleware that limits requests per IP
func (h *Handler) RateLimit(cfg *config.RateLimit) func(http.Handler) http.Handler {
	limiter := newRateLimiter(cfg.RequestsLimit, cfg.WindowSeconds)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !cfg.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			ip := r.RemoteAddr
			// Extract IP without port
			if idx := strings.LastIndex(ip, ":"); idx != -1 {
				ip = ip[:idx]
			}

			if !limiter.allow(ip) {
				RespondTooManyRequests(w, "too many requests, please try again later")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
