package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilyushkaaa/Filmoteka/internal/session/usecase"
	"github.com/ilyushkaaa/Filmoteka/internal/session/usecase/mock"
	usecase2 "github.com/ilyushkaaa/Filmoteka/internal/users/usecase"
	mock2 "github.com/ilyushkaaa/Filmoteka/internal/users/usecase/mock"
	"github.com/ilyushkaaa/Filmoteka/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type fakeHandler struct{}

func (h *fakeHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type errorWriter struct{}

func (w *errorWriter) Header() http.Header {
	return http.Header{}
}

func (w *errorWriter) WriteHeader(_ int) {}

func (w *errorWriter) Write(_ []byte) (int, error) {
	return 0, fmt.Errorf("error")
}

func TestAccessLog(t *testing.T) {
	fakeLogger := zap.NewNop().Sugar()

	middleware := &Middleware{}

	handler := &fakeHandler{}

	req := httptest.NewRequest("GET", "http://films", nil)
	recorder := httptest.NewRecorder()
	middleware.AccessLog(handler).ServeHTTP(recorder, req)
	resp := recorder.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	req = httptest.NewRequest("GET", "http://films", nil)
	ctx := context.WithValue(req.Context(), logger.MyLoggerKey, fakeLogger)
	req = req.WithContext(ctx)
	recorder = httptest.NewRecorder()
	middleware.AccessLog(handler).ServeHTTP(recorder, req)
	resp = recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req = httptest.NewRequest("GET", "http://films", nil)
	middleware.AccessLog(handler).ServeHTTP(&errorWriter{}, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func TestAuth(t *testing.T) {
	fakeLogger := zap.NewNop().Sugar()
	ctrl := gomock.NewController(t)

	sessionUseCase := mock.NewMockSessionUseCase(ctrl)
	userUseCase := mock2.NewMockUserUseCase(ctrl)
	middleware := NewMiddleware(sessionUseCase, userUseCase)

	handler := &fakeHandler{}

	req := httptest.NewRequest("GET", "http://films", nil)
	recorder := httptest.NewRecorder()
	middleware.AuthMiddleware(handler).ServeHTTP(recorder, req)
	resp := recorder.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	req = httptest.NewRequest("GET", "http://films", nil)
	ctx := context.WithValue(req.Context(), logger.MyLoggerKey, fakeLogger)
	req = req.WithContext(ctx)
	recorder = httptest.NewRecorder()
	middleware.AuthMiddleware(handler).ServeHTTP(recorder, req)
	resp = recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	req = httptest.NewRequest("GET", "http://films", nil)
	middleware.AuthMiddleware(handler).ServeHTTP(&errorWriter{}, req)

	sessionUseCase.EXPECT().GetSession("qqq").Return(nil, fmt.Errorf("error"))
	req = httptest.NewRequest("GET", "http://films", nil)
	req.Header = map[string][]string{
		"Cookie": {`session_id="qqq"`},
	}
	ctx = context.WithValue(req.Context(), logger.MyLoggerKey, fakeLogger)
	req = req.WithContext(ctx)
	recorder = httptest.NewRecorder()
	middleware.AuthMiddleware(handler).ServeHTTP(recorder, req)
	resp = recorder.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	sessionUseCase.EXPECT().GetSession("qqq").Return(nil, usecase.ErrNoSession)
	req = httptest.NewRequest("GET", "http://films", nil)
	req.Header = map[string][]string{
		"Cookie": {`session_id="qqq"`},
	}
	ctx = context.WithValue(req.Context(), logger.MyLoggerKey, fakeLogger)
	req = req.WithContext(ctx)
	recorder = httptest.NewRecorder()
	middleware.AuthMiddleware(handler).ServeHTTP(recorder, req)
	resp = recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

}

func TestAdmin(t *testing.T) {
	fakeLogger := zap.NewNop().Sugar()
	ctrl := gomock.NewController(t)

	sessionUseCase := mock.NewMockSessionUseCase(ctrl)
	userUseCase := mock2.NewMockUserUseCase(ctrl)
	middleware := NewMiddleware(sessionUseCase, userUseCase)

	handler := &fakeHandler{}

	req := httptest.NewRequest("GET", "http://films", nil)
	recorder := httptest.NewRecorder()
	middleware.AdminMiddleware(handler).ServeHTTP(recorder, req)
	resp := recorder.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	req = httptest.NewRequest("GET", "http://films", nil)
	ctx := context.WithValue(req.Context(), logger.MyLoggerKey, fakeLogger)
	req = req.WithContext(ctx)
	recorder = httptest.NewRecorder()
	middleware.AdminMiddleware(handler).ServeHTTP(recorder, req)
	resp = recorder.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	userUseCase.EXPECT().GetUserRole(uint64(1)).Return("", fmt.Errorf("error"))
	req = httptest.NewRequest("GET", "http://films", nil)
	ctx = context.WithValue(req.Context(), logger.MyLoggerKey, fakeLogger)
	ctx = context.WithValue(ctx, MyUserKey, uint64(1))
	req = req.WithContext(ctx)
	recorder = httptest.NewRecorder()
	middleware.AdminMiddleware(handler).ServeHTTP(recorder, req)
	resp = recorder.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	userUseCase.EXPECT().GetUserRole(uint64(1)).Return("", usecase2.ErrNoUser)
	req = httptest.NewRequest("GET", "http://films", nil)
	ctx = context.WithValue(req.Context(), logger.MyLoggerKey, fakeLogger)
	ctx = context.WithValue(ctx, MyUserKey, uint64(1))
	req = req.WithContext(ctx)
	recorder = httptest.NewRecorder()
	middleware.AdminMiddleware(handler).ServeHTTP(recorder, req)
	resp = recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	userUseCase.EXPECT().GetUserRole(uint64(1)).Return("default", nil)
	req = httptest.NewRequest("GET", "http://films", nil)
	ctx = context.WithValue(req.Context(), logger.MyLoggerKey, fakeLogger)
	ctx = context.WithValue(ctx, MyUserKey, uint64(1))
	req = req.WithContext(ctx)
	recorder = httptest.NewRecorder()
	middleware.AdminMiddleware(handler).ServeHTTP(recorder, req)
	resp = recorder.Result()
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

}

func TestRequestInit(t *testing.T) {
	middleware := &Middleware{}
	handler := &fakeHandler{}

	req := httptest.NewRequest("GET", "http://films", nil)
	recorder := httptest.NewRecorder()
	middleware.RequestInitMiddleware(handler).ServeHTTP(recorder, req)
	resp := recorder.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
