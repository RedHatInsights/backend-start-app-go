package logging

import (
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewMiddleware(globalLogger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			bytesIn, _ := strconv.Atoi(r.Header.Get("Content-Length"))
			loggerCtx := globalLogger.With().
				Str("remote_ip", r.RemoteAddr).
				Str("url", r.URL.Path).
				Str("method", r.Method)
			contextLogger := loggerCtx.Logger()
			contextLogger.Debug().Msgf("Started %s request %s", r.Method, r.URL.Path)

			t1 := time.Now()
			defer func() {
				duration := time.Since(t1)
				afterLogger := contextLogger.With().
					Dur("latency_ms", duration).
					Int("bytes_in", bytesIn).
					Int("bytes_out", ww.BytesWritten()).
					Logger()

				// prevent the application from exiting
				if rec := recover(); rec != nil {
					afterLogger.Error().
						Bool("panic", true).
						Int("status", http.StatusInternalServerError).
						Msgf("Unhandled panic: %s\n%s", rec, debug.Stack())
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				log.Info().
					Int("status", ww.Status()).
					Msgf("Completed %s request %s in %s with %d",
						r.Method, r.URL.Path, duration.Round(time.Millisecond).String(), ww.Status())
			}()

			ctx := WithLogger(r.Context(), &contextLogger)
			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}
