package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ilyushkaaa/Filmoteka/internal/session/usecase"
	"github.com/ilyushkaaa/Filmoteka/pkg/logger"
	"github.com/ilyushkaaa/Filmoteka/pkg/response"
)

type userKey int
type tokenKey int

const (
	MyUserKey      userKey  = 1
	MySessionIDKey tokenKey = 2
)

func (mw *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zapLogger, err := logger.GetLoggerFromContext(r.Context())
		if err != nil {
			log.Printf("can not get logger from context: %s", err)
			errText := `{"error": "internal error"}`
			err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
			if err != nil {
				log.Printf("can not write response: %s", err)
			}
			return
		}
		sessionCookie, err := r.Cookie("session_id")
		if errors.Is(err, http.ErrNoCookie) {
			zapLogger.Errorf("no cookie in request")
			errText := `{"error": "no cookie in request""}`
			err = response.WriteResponse(w, []byte(errText), http.StatusUnauthorized)
			if err != nil {
				zapLogger.Errorf("can not write response: %s", err)
			}
			return
		}
		if err != nil {
			errText := `{"error": "internal error"}`
			zapLogger.Errorf("error in getting cookie: %s", err)
			err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
			if err != nil {
				zapLogger.Errorf("can not write response: %s", err)
			}
			return
		}
		sessionID := sessionCookie.Value
		mySession, err := mw.sessionUseCase.GetSession(sessionID)
		if errors.Is(err, usecase.ErrNoSession) {
			zapLogger.Errorf("no session for id: %s", sessionID)
			errText := fmt.Sprintf(`{"error": "there is no session for session id %s}`, sessionID)
			err = response.WriteResponse(w, []byte(errText), http.StatusUnauthorized)
			if err != nil {
				zapLogger.Errorf("can not write response: %s", err)
			}
			return
		}
		if err != nil {
			zapLogger.Errorf("error in getting session: %s", err)
			errText := `{"error": "internal error"}`
			err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
			if err != nil {
				zapLogger.Errorf("can not write response: %s", err)
			}
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, MyUserKey, mySession.UserID)
		ctx = context.WithValue(ctx, MySessionIDKey, mySession.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
