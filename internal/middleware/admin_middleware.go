package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/ilyushkaaa/Filmoteka/internal/users/usecase"
	"github.com/ilyushkaaa/Filmoteka/pkg/logger"
	"github.com/ilyushkaaa/Filmoteka/pkg/response"
)

func (mw *Middleware) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		zapLogger, err := logger.GetLoggerFromContext(ctx)
		if err != nil {
			log.Printf("can not get logger from context: %s", err)
			err = response.WriteResponse(w, []byte(`{"error": "internal error"}`), http.StatusInternalServerError)
			if err != nil {
				log.Printf("can not write response: %s", err)
			}
			return
		}
		userID, ok := ctx.Value(MyUserKey).(uint64)
		if !ok {
			zapLogger.Errorf("can not get user id from context")
			err = response.WriteResponse(w, []byte(`{"error": "internal error"}`), http.StatusInternalServerError)
			if err != nil {
				zapLogger.Errorf("can not write response: %s", err)
			}
			return
		}
		role, err := mw.userUseCase.GetUserRole(userID)
		if errors.Is(err, usecase.ErrNoUser) {
			zapLogger.Errorf("user with id %d was not found", userID)
			err = response.WriteResponse(w, []byte(`{"error": "user is not found"}`), http.StatusUnauthorized)
			if err != nil {
				zapLogger.Errorf("can not write response: %s", err)
			}
			return
		}
		if err != nil {
			zapLogger.Errorf("internal error in getting user role: %s", err)
			err = response.WriteResponse(w, []byte(`{"error": "internal error"}`), http.StatusInternalServerError)
			if err != nil {
				zapLogger.Errorf("can not write response: %s", err)
			}
			return
		}
		if role != "admin" {
			zapLogger.Errorf("user is not admin, but wants to use admin resource")
			err = response.WriteResponse(w, []byte(`{"error": "resource is forbidden for you"}`), http.StatusForbidden)
			if err != nil {
				zapLogger.Errorf("can not write response: %s", err)
			}
			return
		}
		next.ServeHTTP(w, r)

	})
}
