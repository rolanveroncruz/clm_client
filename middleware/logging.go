package middleware

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"time"
)

var lumberjackLogger = &lumberjack.Logger{
	Filename:   "logs/log.log", //filename
	MaxSize:    100,            // file size in MB before rotation
	MaxBackups: 10,             // Max number of files kept before being overwritten
	MaxAge:     28,             // Max number of days to keep the files
	Compress:   true,           // Whether to compress log files using gzip
}
var logger = zerolog.New(lumberjackLogger).With().Timestamp().Logger()

// LoggingMiddleware is an HTTP middleware that logs request and response data, including timing and status codes.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		rec := httptest.NewRecorder()
		ctx := r.Context()
		path := r.URL.EscapedPath()
		reqData, _ := httputil.DumpRequest(r, true)

		// log time, path, request_data,
		logger := logger.Log().Timestamp().Str("path", path).Bytes("request_data", reqData)
		defer func(begin time.Time) {
			status := ww.Status()
			tookMs := time.Since(begin).Milliseconds()

			fmt.Printf("%s:%s %d\n", r.Method, path, status)
			// at the end of the call, log took, status_code, and Msg
			logger.Int64("took", tookMs).Int("status_code", status).Msgf("[%d] %s http request for %s took %dms", status, r.Method, path, tookMs)
		}(time.Now())

		// Replace "logger" with a custom type, like ContextKey("logger")
		ctx = context.WithValue(ctx, "logger", logger)
		next.ServeHTTP(rec, r.WithContext(ctx)) // this copies the recorded response to the response writer

		for k, v := range rec.Header() {
			ww.Header()[k] = v
		}
		ww.WriteHeader(rec.Code)
		_, _ = rec.Body.WriteTo(ww)

	})
}
