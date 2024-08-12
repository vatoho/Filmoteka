package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/ilyushkaaa/Filmoteka/pkg/logger"
	"github.com/ilyushkaaa/Filmoteka/pkg/response"
)

func (mw *Middleware) AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zapLogger, err := logger.GetLoggerFromContext(r.Context())
		if err != nil {
			log.Printf("can not get logger from context: %s", err)
			err = response.WriteResponse(w, []byte("internal error"), http.StatusInternalServerError)
			if err != nil {
				log.Printf("can not write response: %s", err)
			}
			return
		}
		zapLogger.Infof("access log middleware start")
		start := time.Now()
		next.ServeHTTP(w, r)
		zapLogger.Infow("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"time", time.Since(start),
		)
	})
}
