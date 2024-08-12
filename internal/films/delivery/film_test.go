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
	"github.com/ilyushkaaa/Filmoteka/internal/films/entity"
	"github.com/ilyushkaaa/Filmoteka/internal/films/usecase"
	"github.com/ilyushkaaa/Filmoteka/internal/films/usecase/mock"
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

func TestGetFilms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockFilmUseCase(ctrl)
	testHandler := NewFilmHandler(testUseCase)

	// can not read request body

	request := httptest.NewRequest(http.MethodGet, "/films", &errorReader{})
	respWriter := httptest.NewRecorder()
	testHandler.GetFilms(respWriter, request)
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

	request = httptest.NewRequest(http.MethodGet, "/films", &errorReader{})
	respErr := &errorWriter{}
	testHandler.GetFilms(respErr, request)
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

	testUseCase.EXPECT().GetFilms("").Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/films", nil)
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilms(respWriter, request.WithContext(ctx))
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

	testUseCase.EXPECT().GetFilms("").Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/films", nil)
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	testHandler.GetFilms(respErr, request.WithContext(ctx))
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

	films := make([]entity.Film, 0)
	testUseCase.EXPECT().GetFilms("").Return(films, nil)
	request = httptest.NewRequest(http.MethodGet, "/films", nil)
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilms(respWriter, request.WithContext(ctx))
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

func TestGetFilmByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockFilmUseCase(ctrl)
	testHandler := NewFilmHandler(testUseCase)

	// can not read request body

	request := httptest.NewRequest(http.MethodGet, "/film/1", &errorReader{})
	respWriter := httptest.NewRecorder()
	testHandler.GetFilmByID(respWriter, request)
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

	request = httptest.NewRequest(http.MethodGet, "/film/1", &errorReader{})
	respErr := &errorWriter{}
	testHandler.GetFilmByID(respErr, request)
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

	request = httptest.NewRequest(http.MethodGet, "/film/bad_id", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "bad_id"})
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilmByID(respWriter, request.WithContext(ctx))
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
	testUseCase.EXPECT().GetFilmByID(id).Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/film/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilmByID(respWriter, request.WithContext(ctx))
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

	testUseCase.EXPECT().GetFilmByID(id).Return(nil, usecase.ErrFilmNotFound)
	request = httptest.NewRequest(http.MethodGet, "/film/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilmByID(respWriter, request.WithContext(ctx))
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

	film := &entity.Film{}
	testUseCase.EXPECT().GetFilmByID(id).Return(film, nil)
	request = httptest.NewRequest(http.MethodGet, "/film/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilmByID(respWriter, request.WithContext(ctx))
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

func TestAddFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockFilmUseCase(ctrl)
	testHandler := NewFilmHandler(testUseCase)

	request := httptest.NewRequest(http.MethodPost, "/film", &errorReader{})
	respWriter := httptest.NewRecorder()
	testHandler.AddFilm(respWriter, request)
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

	request = httptest.NewRequest(http.MethodPost, "/film", &errorReader{})
	respErr := &errorWriter{}
	testHandler.AddFilm(respErr, request)
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

	request = httptest.NewRequest(http.MethodPost, "/film", &errorReader{})
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddFilm(respWriter, request.WithContext(ctx))
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

	request = httptest.NewRequest(http.MethodPost, "/film", strings.NewReader(`{"`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddFilm(respWriter, request.WithContext(ctx))
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

	request = httptest.NewRequest(http.MethodPost, "/film", strings.NewReader(`{"name":"ded"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddFilm(respWriter, request.WithContext(ctx))
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

	film := entity.Film{
		Name:          "qqq",
		Description:   "fff",
		DateOfRelease: time.Time{}.Add(time.Hour),
		Rating:        5.1,
	}
	testUseCase.EXPECT().AddFilm(film, []uint64{}).Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodPost, "/film", strings.NewReader(
		`{"name":"qqq","description":"fff","date_of_release":"0001-01-01T01:00:00Z","rating":5.1}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddFilm(respWriter, request.WithContext(ctx))
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

	testUseCase.EXPECT().AddFilm(film, []uint64{}).Return(nil, usecase.ErrBadFilmAddData)
	request = httptest.NewRequest(http.MethodPost, "/film", strings.NewReader(
		`{"name":"qqq","description":"fff","date_of_release":"0001-01-01T01:00:00Z","rating":5.1}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddFilm(respWriter, request.WithContext(ctx))
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

	filmAdded := film
	filmAdded.ID = 1
	testUseCase.EXPECT().AddFilm(film, []uint64{}).Return(&filmAdded, nil)
	request = httptest.NewRequest(http.MethodPost, "/film", strings.NewReader(
		`{"name":"qqq","description":"fff","date_of_release":"0001-01-01T01:00:00Z","rating":5.1}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddFilm(respWriter, request.WithContext(ctx))
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

func TestUpdateFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockFilmUseCase(ctrl)
	testHandler := NewFilmHandler(testUseCase)

	request := httptest.NewRequest(http.MethodPut, "/film", &errorReader{})
	respWriter := httptest.NewRecorder()
	testHandler.UpdateFilm(respWriter, request)
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

	request = httptest.NewRequest(http.MethodPut, "/film", &errorReader{})
	respErr := &errorWriter{}
	testHandler.UpdateFilm(respErr, request)
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

	request = httptest.NewRequest(http.MethodPut, "/film", &errorReader{})
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.UpdateFilm(respWriter, request.WithContext(ctx))
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

	request = httptest.NewRequest(http.MethodPut, "/film", strings.NewReader(`{"`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.UpdateFilm(respWriter, request.WithContext(ctx))
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

	request = httptest.NewRequest(http.MethodPut, "/film", strings.NewReader(`{"name":"ded"}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.UpdateFilm(respWriter, request.WithContext(ctx))
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

	film := entity.Film{
		ID:            1,
		Name:          "qqq",
		Description:   "fff",
		DateOfRelease: time.Time{}.Add(time.Hour),
		Rating:        5.1,
	}
	testUseCase.EXPECT().UpdateFilm(film, []uint64{}).Return(fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodPut, "/film", strings.NewReader(
		`{"id":1,"name":"qqq","description":"fff","date_of_release":"0001-01-01T01:00:00Z","rating":5.1}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.UpdateFilm(respWriter, request.WithContext(ctx))
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

	testUseCase.EXPECT().UpdateFilm(film, []uint64{}).Return(usecase.ErrBadFilmUpdateData)
	request = httptest.NewRequest(http.MethodPut, "/film", strings.NewReader(
		`{"id":1,"name":"qqq","description":"fff","date_of_release":"0001-01-01T01:00:00Z","rating":5.1}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.UpdateFilm(respWriter, request.WithContext(ctx))
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

	testUseCase.EXPECT().UpdateFilm(film, []uint64{}).Return(nil)
	request = httptest.NewRequest(http.MethodPut, "/film", strings.NewReader(
		`{"id":1,"name":"qqq","description":"fff","date_of_release":"0001-01-01T01:00:00Z","rating":5.1}`))
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.UpdateFilm(respWriter, request.WithContext(ctx))
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

func TestGetFilmsBySearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockFilmUseCase(ctrl)
	testHandler := NewFilmHandler(testUseCase)

	request := httptest.NewRequest(http.MethodGet, "/films/search", &errorReader{})
	respWriter := httptest.NewRecorder()
	testHandler.GetFilmsBySearch(respWriter, request)
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

	request = httptest.NewRequest(http.MethodGet, "/films/search", &errorReader{})
	respErr := &errorWriter{}
	testHandler.GetFilmsBySearch(respErr, request)
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

	testUseCase.EXPECT().GetFilmsBySearch("").Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/films/search", nil)
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilmsBySearch(respWriter, request.WithContext(ctx))
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

	testUseCase.EXPECT().GetFilmsBySearch("").Return(nil, usecase.ErrFilmsNotFound)
	request = httptest.NewRequest(http.MethodGet, "/films/search", nil)
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilmsBySearch(respWriter, request.WithContext(ctx))
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

	testUseCase.EXPECT().GetFilmsBySearch("").Return([]entity.Film{}, nil)
	request = httptest.NewRequest(http.MethodGet, "/films/search", nil)
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilmsBySearch(respWriter, request.WithContext(ctx))
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

func TestDeleteFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := mock.NewMockFilmUseCase(ctrl)
	testHandler := NewFilmHandler(testUseCase)

	request := httptest.NewRequest(http.MethodDelete, "/film/1", &errorReader{})
	respWriter := httptest.NewRecorder()
	testHandler.DeleteFilm(respWriter, request)
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

	request = httptest.NewRequest(http.MethodDelete, "/films/1", &errorReader{})
	respErr := &errorWriter{}
	testHandler.DeleteFilm(respErr, request)
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

	request = httptest.NewRequest(http.MethodDelete, "/film/bad_id", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "bad_id"})
	ctx := request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.DeleteFilm(respWriter, request.WithContext(ctx))
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
	testUseCase.EXPECT().DeleteFilm(id).Return(fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/film/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, logger2.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.DeleteFilm(respWriter, request.WithContext(ctx))
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
