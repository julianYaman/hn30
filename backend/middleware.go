package main

import (
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Create a map to store rate limiters for each visitor
var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, 0}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytesWritten += n
	return n, err
}

// Unified HTTP logging middleware with structured logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)

		// Extract real IP
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
			"event_type", "http_request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_ip", ip,
		)

		logger.Info("request received",
			"event", "request_started",
			"user_agent", r.UserAgent(),
			"referer", r.Referer(),
		)

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		statusCode := lrw.statusCode

		logLevel := slog.LevelInfo
		if statusCode >= 500 {
			logLevel = slog.LevelError
		} else if statusCode >= 400 {
			logLevel = slog.LevelWarn
		}

		logger.Log(r.Context(), logLevel, "request completed",
			"event", "request_completed",
			"status_code", statusCode,
			"bytes_written", lrw.bytesWritten,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

// Middleware to limit requests from a single IP
func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the real IP address from the X-Forwarded-For header.
		// This is crucial when running behind a reverse proxy like Traefik.
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			// Fallback to RemoteAddr for local development or direct connections.
			ip = r.RemoteAddr
		}

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
			"event_type", "rate_limit",
			"remote_ip", ip,
			"path", r.URL.Path,
		)

		mu.Lock()
		// Check if the visitor has a limiter yet
		newVisitor := false
		if _, found := visitors[ip]; !found {
			// Allow 10 requests per minute (1 request every 6 seconds)
			visitors[ip] = rate.NewLimiter(rate.Every(6*time.Second), 1)
			newVisitor = true
		}
		limiter := visitors[ip]
		mu.Unlock()

		if newVisitor {
			logger.Info("new rate limiter created",
				"event", "rate_limiter_created",
				"rate", "1 request per 6 seconds",
				"burst", 1,
			)
		}

		// If the visitor's limiter doesn't allow the request, block them.
		if !limiter.Allow() {
			logger.Warn("rate limit exceeded",
				"event", "rate_limit_exceeded",
				"tokens_available", limiter.Tokens(),
			)
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		logger.Info("rate limit check passed",
			"event", "rate_limit_allowed",
			"tokens_remaining", limiter.Tokens(),
		)

		// Otherwise, serve the request
		next.ServeHTTP(w, r)
	})
}
