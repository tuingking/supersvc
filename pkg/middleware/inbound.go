package middleware

import (
	"context"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tuingking/supersvc/pkg/ctxkey"
)

// NOTE: this package using zerolog package
// set "x-request-id" to context and zerolog
func Inbound(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestId := r.Header.Get(ctxkey.XRequestID.String())
		if requestId == "" {
			requestId = uuid.New().String()
		}

		// create sub logger with requestId
		subLogger := log.With().Str(ctxkey.XRequestID.String(), requestId).Logger()

		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			subLogger.Err(err).Msg("err: httputil.DumpRequest")
		}

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		tStart := time.Now()
		defer func() {
			tEnd := time.Now()

			// recover and record stack traces in case of a panic
			if rec := recover(); rec != nil {
				subLogger.Error().
					Str("type", "error").
					Interface("recover_info", rec).
					Bytes("debug_stack", debug.Stack()).
					Msg("log system error")
				http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			// log end request
			subLogger.Info().
				Fields(map[string]interface{}{
					"remote_ip":  r.RemoteAddr,
					"url":        r.URL.Path,
					"proto":      r.Proto,
					"method":     r.Method,
					"user_agent": r.Header.Get(ctxkey.UserAgent.String()),
					"status":     ww.Status(),
					"latency_ms": float64(tEnd.Sub(tStart).Nanoseconds()) / 1000000.0,
					"bytes_in":   r.Header.Get(ctxkey.ContentLength.String()),
					"bytes_out":  ww.BytesWritten(),
					"dump":       string(dump),
				}).
				Msg("http inbound")
		}()

		// pass sub-logger & request id by context
		ctx = context.WithValue(ctx, ctxkey.XRequestID, requestId)
		ctx = context.WithValue(ctx, ctxkey.ZeroLogSubLogger, subLogger)
		// ctx = context.WithValue(ctx, ctxkey.ZeroLogSubLoggerCtx, subLogger.WithContext(ctx)) // alternative #2 (send context instead of zerlog.Logger)

		next.ServeHTTP(ww, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
