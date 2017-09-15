package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
)

type key int

const loggerIDKey key = 82

// Logger middleware.
type Logger struct {
	h      http.Handler
	logger log.Logger
}

//SetLogger ...
func (l *Logger) SetLogger(logger log.Logger) {
	l.logger = logger
}

// wrapper to capture status.
type wrapper struct {
	http.ResponseWriter
	written int
	status  int
}

// capture status.
func (w *wrapper) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// capture written bytes.
func (w *wrapper) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.written += n
	return n, err
}

// NewLogger middleware with the given log.Logger.
func NewLogger(logger log.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &Logger{
			logger: logger,
			h:      h,
		}
	}
}

// ServeHTTP implementation.
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	res := &wrapper{w, 0, 200}

	// get the context since we'll use it a few times
	ctx := r.Context()

	logger := log.With(l.logger)

	// log the request.start
	logger.Log(
		"method", r.Method,
		"uri", r.RequestURI,
		//"ip", realip.RealIpFromContext(ctx),
		"evt", "request.start",
	)

	// continue to the next middleware
	ctx = context.WithValue(ctx, loggerIDKey, logger)
	l.h.ServeHTTP(res, r.WithContext(ctx))

	// log the request.end
	logger.Log(
		"status", res.status,
		"size", res.written,
		"duration", time.Since(start),
		"evt", "request.end",
	)
}

// LoggerFromRequest can be used to obtain the Log from the request.
func LoggerFromRequest(r *http.Request) log.Logger {
	return r.Context().Value(loggerIDKey).(log.Logger)
}

// LoggerFromContext can be used to obtain the Log from the request.
func LoggerFromContext(ctx context.Context) log.Logger {
	return ctx.Value(loggerIDKey).(log.Logger)
}
