package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	cyan   = "\033[36m"
	reset  = "\033[0m"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func prettyJSON(data []byte) string {
	var out bytes.Buffer
	if json.Valid(data) {
		if err := json.Indent(&out, data, "", "  "); err == nil {
			return out.String()
		}
	}
	return string(data)
}

func maybePretty(body []byte, contentType string) string {
	if strings.Contains(contentType, "application/json") {
		return prettyJSON(body)
	}
	return string(body)
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 400 && code < 500:
		return yellow
	case code >= 500:
		return red
	default:
		return cyan
	}
}

// generateRequestID creates a timestamp-based ID
func generateRequestID() string {
	return fmt.Sprintf("req-%d", time.Now().Unix())
}

func requestLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			reqID := generateRequestID()

			var reqBodyRaw []byte
			if r.Body != nil {
				reqBodyRaw, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewReader(reqBodyRaw))
			}

			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           &bytes.Buffer{},
			}

			next.ServeHTTP(rw, r)

			duration := time.Since(start).Milliseconds()
			color := colorForStatus(rw.statusCode)

			log.Printf("%s---[Request %s Start]----------------------------%s", color, reqID, reset)
			log.Printf("%s[%s] %s %s%s", color, reqID, r.Method, r.URL.Path, reset)
			log.Printf("%sHeaders:%s %v", color, reset, r.Header)
			log.Printf(
				"%sRequest Body:%s\n%s",
				color,
				reset,
				maybePretty(reqBodyRaw, r.Header.Get("Content-Type")),
			)
			log.Printf(
				"%sResponse (%d):%s\n%s",
				color,
				rw.statusCode,
				reset,
				maybePretty(rw.body.Bytes(), rw.Header().Get("Content-Type")),
			)
			log.Printf("%sDuration:%s %d ms", color, reset, duration)
			log.Printf(
				"%s---[Request %s End]------------------------------%s\n",
				color,
				reqID,
				reset,
			)
		})
	}
}
