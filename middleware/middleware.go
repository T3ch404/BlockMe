package middleware

import (
	"BlockMe/config"
	"BlockMe/to"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

func BlockCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		theirIpStr := IpFromContext(r.Context()).String()
		if theirIpStr == "" {
			fmt.Println("IP not found in context")
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}

		// Skip block check on reset requests
		if strings.Contains(r.URL.String(), "reset") {
			fmt.Printf("reset link, skipping block check\n")
			next.ServeHTTP(w, r)
		}

		var exists bool
		sqlStatement := `SELECT EXISTS (SELECT 1 FROM ip_blocklist WHERE ip = $1)`
		err := to.DB.QueryRow(sqlStatement, theirIpStr).Scan(&exists)
		if err == nil && exists {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RealIpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse X-Forwarded-For header to net.IP
		theirIPStr := r.Header.Get("X-Forwarded-For")
		if theirIPStr == "" {
			var err error
			theirIPStr, _, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				fmt.Println("Error parsing remote address")
				http.Error(w, "Error", http.StatusInternalServerError)
				return
			}
		}

		// Verify IP string is an actual IP
		theirIP := net.ParseIP(theirIPStr)
		if theirIP == nil {
			http.Error(w, "You didn't pass a valid IP in the X-Forwarded-For header. Stop doing that.", http.StatusBadRequest)
			return
		}

		// Add their IP to the request context
		ctx := contextWithIP(r.Context(), theirIP)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *CustomResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.StatusCode = code
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	// Only set StatusCode to 200 if it was not already set
	if w.StatusCode == 0 {
		w.StatusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customWriter := &CustomResponseWriter{ResponseWriter: w}

		start := time.Now()

		next.ServeHTTP(customWriter, r)

		theirIpStr := IpFromContext(r.Context()).String()
		if theirIpStr == "" {
			fmt.Println("IP not found in context")
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}

		duration := time.Since(start)

		if config.Env.LogFormat == "text_pretty" || config.Env.LogLocation == "file" {
			// Simple human-readable logging when LOG_FORMAT is text_pretty or slog is going to file
			fmt.Printf("%s\t%s\t%s\t%d\t%s\n", theirIpStr, r.Method, r.URL.String(), customWriter.StatusCode, duration.String())
		}

		if config.Env.LogFormat != "text_pretty" {
			slog.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"request",
				slog.String("source", theirIpStr),
				slog.String("method", r.Method),
				slog.String("host", r.Host),
				slog.String("path", r.URL.String()),
				slog.Int("response_code", customWriter.StatusCode),
				slog.String("start", start.Format(time.RFC3339)),
				slog.String("duration", duration.String()),
				slog.String("user_agent", r.UserAgent()),
				slog.Int("bytes", int(r.ContentLength)),
			)
		}
	})
}

type key int

const ipKey key = iota

func contextWithIP(ctx context.Context, ip net.IP) context.Context {
	return context.WithValue(ctx, ipKey, ip)
}

func IpFromContext(ctx context.Context) net.IP {
	return ctx.Value(ipKey).(net.IP)
}
