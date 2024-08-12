package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ilyushkaaa/Filmoteka/pkg/logger"
	"github.com/ilyushkaaa/Filmoteka/pkg/response"
	"go.uber.org/zap"
)

func (mw *Middleware) RequestInitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		myLogger, err := logger.InitLogger()
		if err != nil {
			err = response.WriteResponse(w, []byte("internal error"), http.StatusInternalServerError)
			if err != nil {
				log.Printf("can not write response: %s", err)
			}
			return
		}
		requestID := uuid.New().String()
		myLogger = myLogger.With(zap.String("request-id", requestID))
		ctx := r.Context()
		ctx = context.WithValue(ctx, logger.MyLoggerKey, myLogger)
		myLogger.Infof("request init middleware call")
		next.ServeHTTP(w, r.WithContext(ctx))
		loggerErr := myLogger.Sync()
		if loggerErr != nil {
			log.Println("error in logger sync")
		}
	})
}
