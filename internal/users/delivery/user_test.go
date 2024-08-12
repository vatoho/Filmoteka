package delivery

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	usecase2 "github.com/ilyushkaaa/Filmoteka/internal/session/usecase/mock"
	"github.com/ilyushkaaa/Filmoteka/internal/users/entity"
	"github.com/ilyushkaaa/Filmoteka/internal/users/usecase"
	"github.com/ilyushkaaa/Filmoteka/internal/users/usecase/mock"
	logger2 "github.com/ilyushkaaa/Filmoteka/pkg/logger"
	"go.uber.org/zap"
)

type errorReader struct{}

func (er *errorReader) Read(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

type errorWriter struct{}

func (w *errorWriter) Header() http.Header {
	return http.Header{}
}

func (w *errorWriter) WriteHeader(_ int) {}

func (w *errorWriter) Write(_ []byte) (int, error) {
	return 0, fmt.Errorf("error")
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockUserUseCase(ctrl)
	sessionTestUseCase := usecase2.NewMockSessionUseCase(ctrl)
	testHandler := NewUserHandler(testUseCase, sessionTestUseCase)

	// can not read request body
	request := httptest.NewRequest(http.MethodPost, "/login", &errorReader{})
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter := httptest.NewRecorder()
	testHandler.Login(respWriter, request.WithContext(ctx))
	resp := respWriter.Result()
	_, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}

	request = httptest.NewRequest(http.MethodPost, "/login", &errorReader{})
	respWriter = httptest.NewRecorder()
	testHandler.Login(respWriter, request)
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}

	request = httptest.NewRequest(http.MethodPost, "/login", &errorReader{})
	respErr := &errorWriter{}
	testHandler.Login(respErr, request)
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}

	request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("{"))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.Login(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 401 {
		t.Errorf("expected status %d, got status %d", http.StatusUnauthorized, resp.StatusCode)
	}

	request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("{"))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	testHandler.Login(respErr, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 401 {
		t.Errorf("expected status %d, got status %d", http.StatusUnauthorized, resp.StatusCode)
	}

	request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"dwdwdwd"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	testHandler.Login(respErr, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 401 {
		t.Errorf("expected status %d, got status %d", http.StatusUnauthorized, resp.StatusCode)
	}

	testUseCase.EXPECT().Login("hello12", "qqqqqqqqq").Return(nil, usecase.ErrBadCredentials)
	request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"hello12","password":"qqqqqqqqq"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.Login(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 401 {
		t.Errorf("expected status %d, got status %d", http.StatusUnauthorized, resp.StatusCode)
	}

	testUseCase.EXPECT().Login("hello12", "qqqqqqqqq").Return(nil, fmt.Errorf("internal server error"))
	request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"hello12","password":"qqqqqqqqq"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.Login(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}

	testUseCase.EXPECT().Login("hello12", "qqqqqqqqq").Return(nil, fmt.Errorf("internal server error"))
	request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"hello12","password":"qqqqqqqqq"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	testHandler.Login(respErr, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}

	loggedInUser := &entity.User{
		ID:       1,
		Username: "some_username",
	}
	testUseCase.EXPECT().Login("some_username", "aaaaaaaa").Return(loggedInUser, nil)
	sessionTestUseCase.EXPECT().CreateSession(loggedInUser.ID).Return("", fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"some_username", "password":"aaaaaaaa"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.Login(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}

	testUseCase.EXPECT().Login("some_username", "aaaaaaaa").Return(loggedInUser, nil)
	sessionTestUseCase.EXPECT().CreateSession(loggedInUser.ID).Return("some_token", nil)
	request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"some_username", "password":"aaaaaaaa"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.Login(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status %d, got status %d", http.StatusOK, resp.StatusCode)
	}
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockUserUseCase(ctrl)
	sessionTestUseCase := usecase2.NewMockSessionUseCase(ctrl)
	testHandler := NewUserHandler(testUseCase, sessionTestUseCase)

	// can not read request body
	request := httptest.NewRequest(http.MethodPost, "/register", &errorReader{})
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter := httptest.NewRecorder()
	testHandler.Register(respWriter, request.WithContext(ctx))
	resp := respWriter.Result()
	_, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}

	request = httptest.NewRequest(http.MethodPost, "/register", &errorReader{})
	respWriter = httptest.NewRecorder()
	testHandler.Register(respWriter, request)
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}

	request = httptest.NewRequest(http.MethodPost, "/register", &errorReader{})
	respErr := &errorWriter{}
	testHandler.Register(respErr, request)
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}

	request = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("{"))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.Register(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 401 {
		t.Errorf("expected status %d, got status %d", http.StatusUnauthorized, resp.StatusCode)
	}

	testUseCase.EXPECT().Register("hello12", "qqqqqqqqq").Return(nil, usecase.ErrUserAlreadyExists)
	request = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"hello12","password":"qqqqqqqqq"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.Register(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 422 {
		t.Errorf("expected status %d, got status %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}

	testUseCase.EXPECT().Register("hello12", "qqqqqqqqq").Return(nil, fmt.Errorf("internal server error"))
	request = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"hello12","password":"qqqqqqqqq"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.Register(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}
	loggedInUser := &entity.User{
		ID:       1,
		Username: "some_username",
	}
	testUseCase.EXPECT().Register("some_username", "aaaaaaaa").Return(loggedInUser, nil)
	sessionTestUseCase.EXPECT().CreateSession(loggedInUser.ID).Return("", fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"some_username", "password":"aaaaaaaa"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.Register(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockUserUseCase(ctrl)
	sessionTestUseCase := usecase2.NewMockSessionUseCase(ctrl)
	testHandler := NewUserHandler(testUseCase, sessionTestUseCase)

	// can not read request body

	request := httptest.NewRequest(http.MethodPost, "/logout", &errorReader{})
	respWriter := httptest.NewRecorder()
	testHandler.Logout(respWriter, request)
	resp := respWriter.Result()
	_, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}

	request = httptest.NewRequest(http.MethodPost, "/logout", &errorReader{})
	respErr := &errorWriter{}
	testHandler.Logout(respErr, request)
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}

	request = httptest.NewRequest(http.MethodPost, "/logout", &errorReader{})
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.Logout(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
	}

}
