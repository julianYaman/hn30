package utils

import (
	"log"
	"net/http"
	"time"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

func LogInfo(format string, v ...interface{}) {
	log.Printf(ColorGreen+"[INFO]"+ColorReset+" "+format, v...)
}

func LogWarn(format string, v ...interface{}) {
	log.Printf(ColorYellow+"[WARN]"+ColorReset+" "+format, v...)
}

func LogError(format string, v ...interface{}) {
	log.Printf(ColorRed+"[ERROR]"+ColorReset+" "+format, v...)
}

func LogComponent(component, format string, v ...interface{}) {
	log.Printf(ColorCyan+"[%s]"+ColorReset+" "+format, append([]interface{}{component}, v...)...)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		statusCode := lrw.statusCode
		statusColor := ColorGreen
		if statusCode >= 400 && statusCode < 500 {
			statusColor = ColorYellow
		} else if statusCode >= 500 {
			statusColor = ColorRed
		}

		LogComponent("HTTP", "%s %s %s%d%s %v", r.Method, r.URL.Path, statusColor, statusCode, ColorReset, duration)
	})
}
