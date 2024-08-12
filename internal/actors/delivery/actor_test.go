package delivery

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/ilyushkaaa/Filmoteka/internal/actors/entity"
	"github.com/ilyushkaaa/Filmoteka/internal/actors/usecase"
	"github.com/ilyushkaaa/Filmoteka/internal/actors/usecase/mock"
	"github.com/ilyushkaaa/Filmoteka/internal/dto"
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

func TestGetActors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockActorUseCase(ctrl)
	testHandler := NewActorHandler(testUseCase)

	request := httptest.NewRequest(http.MethodGet, "/actors", &errorReader{})
	respWriter := httptest.NewRecorder()
	testHandler.GetActors(respWriter, request)
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

	request = httptest.NewRequest(http.MethodGet, "/actors", &errorReader{})
	respErr := &errorWriter{}
	testHandler.GetActors(respErr, request)
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

	testUseCase.EXPECT().GetActors().Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/actors", nil)
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetActors(respWriter, request.WithContext(ctx))
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
		t.Errorf("expected status %d, got status %d", http.StatusUnauthorized, resp.StatusCode)
	}

	testUseCase.EXPECT().GetActors().Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/actors", nil)
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	testHandler.GetActors(respErr, request.WithContext(ctx))
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

	actors := make([]dto.ActorWithFilms, 0)
	testUseCase.EXPECT().GetActors().Return(actors, nil)
	request = httptest.NewRequest(http.MethodGet, "/actors", nil)
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetActors(respWriter, request.WithContext(ctx))
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

func TestGetActorByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockActorUseCase(ctrl)
	testHandler := NewActorHandler(testUseCase)

	request := httptest.NewRequest(http.MethodGet, "/actor/1", &errorReader{})
	respWriter := httptest.NewRecorder()
	testHandler.GetActorByID(respWriter, request)
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

	request = httptest.NewRequest(http.MethodGet, "/actor/1", &errorReader{})
	respErr := &errorWriter{}
	testHandler.GetActorByID(respErr, request)
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

	request = httptest.NewRequest(http.MethodGet, "/actor/bad_id", nil)
	request = mux.SetURLVars(request, map[string]string{"ACTOR_ID": "bad_id"})
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetActorByID(respWriter, request.WithContext(ctx))
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
	if resp.StatusCode != 400 {
		t.Errorf("expected status %d, got status %d", http.StatusBadRequest, resp.StatusCode)
	}

	var id uint64 = 1
	testUseCase.EXPECT().GetActorByID(id).Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/actor/1", nil)
	request = mux.SetURLVars(request, map[string]string{"ACTOR_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetActorByID(respWriter, request.WithContext(ctx))
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

	testUseCase.EXPECT().GetActorByID(id).Return(nil, usecase.ErrActorNotFound)
	request = httptest.NewRequest(http.MethodGet, "/actor/1", nil)
	request = mux.SetURLVars(request, map[string]string{"ACTOR_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetActorByID(respWriter, request.WithContext(ctx))
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
	if resp.StatusCode != 404 {
		t.Errorf("expected status %d, got status %d", http.StatusNotFound, resp.StatusCode)
	}

	actor := dto.ActorWithFilms{}
	testUseCase.EXPECT().GetActorByID(id).Return(&actor, nil)
	request = httptest.NewRequest(http.MethodGet, "/actor/1", nil)
	request = mux.SetURLVars(request, map[string]string{"ACTOR_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetActorByID(respWriter, request.WithContext(ctx))
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

func TestAddActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockActorUseCase(ctrl)
	testHandler := NewActorHandler(testUseCase)

	request := httptest.NewRequest(http.MethodPost, "/actor", &errorReader{})
	respWriter := httptest.NewRecorder()
	testHandler.AddActor(respWriter, request)
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

	request = httptest.NewRequest(http.MethodPost, "/actor", &errorReader{})
	respErr := &errorWriter{}
	testHandler.AddActor(respErr, request)
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

	request = httptest.NewRequest(http.MethodPost, "/actor", &errorReader{})
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddActor(respWriter, request.WithContext(ctx))
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

	request = httptest.NewRequest(http.MethodPost, "/actor", strings.NewReader(`{"`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddActor(respWriter, request.WithContext(ctx))
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
	if resp.StatusCode != 400 {
		t.Errorf("expected status %d, got status %d", http.StatusBadRequest, resp.StatusCode)
	}

	request = httptest.NewRequest(http.MethodPost, "/actor", strings.NewReader(`{"name":"ded"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddActor(respWriter, request.WithContext(ctx))
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

	actor := entity.Actor{
		Name:     "Aaa",
		Surname:  "Aaa",
		Gender:   "male",
		Birthday: time.Time{}.Add(time.Hour),
	}
	testUseCase.EXPECT().AddActor(actor).Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodPost, "/actor", strings.NewReader(
		`{"name":"Aaa","surname":"Aaa","birthday":"0001-01-01T01:00:00Z","gender":"male"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddActor(respWriter, request.WithContext(ctx))
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

	actorAdded := actor
	actorAdded.ID = 1
	testUseCase.EXPECT().AddActor(actor).Return(&actorAdded, nil)
	request = httptest.NewRequest(http.MethodPost, "/actor", strings.NewReader(
		`{"name":"Aaa","surname":"Aaa","birthday":"0001-01-01T01:00:00Z","gender":"male"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddActor(respWriter, request.WithContext(ctx))
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

func TestUpdateActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockActorUseCase(ctrl)
	testHandler := NewActorHandler(testUseCase)

	request := httptest.NewRequest(http.MethodPut, "/actor", &errorReader{})
	respWriter := httptest.NewRecorder()
	testHandler.UpdateActor(respWriter, request)
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

	request = httptest.NewRequest(http.MethodPut, "/actor", &errorReader{})
	respErr := &errorWriter{}
	testHandler.UpdateActor(respErr, request)
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

	request = httptest.NewRequest(http.MethodPut, "/actor", &errorReader{})
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.UpdateActor(respWriter, request.WithContext(ctx))
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

	request = httptest.NewRequest(http.MethodPut, "/actor", strings.NewReader(`{"`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.UpdateActor(respWriter, request.WithContext(ctx))
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
	if resp.StatusCode != 400 {
		t.Errorf("expected status %d, got status %d", http.StatusBadRequest, resp.StatusCode)
	}

	request = httptest.NewRequest(http.MethodPut, "/actor", strings.NewReader(`{"name":"ded"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.UpdateActor(respWriter, request.WithContext(ctx))
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

	actor := entity.Actor{
		ID:       1,
		Name:     "Aaa",
		Surname:  "Aaa",
		Gender:   "male",
		Birthday: time.Time{}.Add(time.Hour),
	}
	testUseCase.EXPECT().UpdateActor(actor).Return(fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodPut, "/actor", strings.NewReader(
		`{"id":1,"name":"Aaa","surname":"Aaa","birthday":"0001-01-01T01:00:00Z","gender":"male"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.UpdateActor(respWriter, request.WithContext(ctx))
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

	testUseCase.EXPECT().UpdateActor(actor).Return(nil)
	request = httptest.NewRequest(http.MethodPut, "/actor", strings.NewReader(
		`{"id":1,"name":"Aaa","surname":"Aaa","birthday":"0001-01-01T01:00:00Z","gender":"male"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.UpdateActor(respWriter, request.WithContext(ctx))
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

func TestDeleteActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockActorUseCase(ctrl)
	testHandler := NewActorHandler(testUseCase)

	request := httptest.NewRequest(http.MethodDelete, "/actor/1", &errorReader{})
	respWriter := httptest.NewRecorder()
	testHandler.DeleteActor(respWriter, request)
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

	request = httptest.NewRequest(http.MethodDelete, "/actor/1", &errorReader{})
	respErr := &errorWriter{}
	testHandler.DeleteActor(respErr, request)
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

	request = httptest.NewRequest(http.MethodDelete, "/actor/bad_id", nil)
	request = mux.SetURLVars(request, map[string]string{"ACTOR_ID": "bad_id"})
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.DeleteActor(respWriter, request.WithContext(ctx))
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
	if resp.StatusCode != 400 {
		t.Errorf("expected status %d, got status %d", http.StatusBadRequest, resp.StatusCode)
	}

	var id uint64 = 1
	testUseCase.EXPECT().DeleteActor(id).Return(fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/actor/1", nil)
	request = mux.SetURLVars(request, map[string]string{"ACTOR_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.DeleteActor(respWriter, request.WithContext(ctx))
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

	testUseCase.EXPECT().DeleteActor(id).Return(nil)
	request = httptest.NewRequest(http.MethodGet, "/actor/1", nil)
	request = mux.SetURLVars(request, map[string]string{"ACTOR_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.DeleteActor(respWriter, request.WithContext(ctx))
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
