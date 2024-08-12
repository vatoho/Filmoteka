package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/ilyushkaaa/Filmoteka/internal/dto"
	_ "github.com/ilyushkaaa/Filmoteka/internal/films/entity"
	"github.com/ilyushkaaa/Filmoteka/internal/films/usecase"
	"github.com/ilyushkaaa/Filmoteka/pkg/logger"
	"github.com/ilyushkaaa/Filmoteka/pkg/response"

	"github.com/gorilla/mux"
)

type FilmHandler struct {
	filmUseCase usecase.FilmUseCase
}

func NewFilmHandler(filmUseCase usecase.FilmUseCase) *FilmHandler {
	return &FilmHandler{
		filmUseCase: filmUseCase,
	}
}

// GetFilms @Summary Получить все фильмы
// @Description Получить список всех фильмов
// @Tags films
// @Accept json
// @Produce json
// @Success 200 {array} entity.Film
// @Failure 400 {object} string "Передан неверный параметр сортировки"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/films [get]
func (h *FilmHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	zapLogger, err := logger.GetLoggerFromContext(r.Context())
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		err = response.WriteResponse(w, []byte(`"error":"internal error"`), http.StatusInternalServerError)
		if err != nil {
			log.Printf("can not write response: %s", err)
		}
		return
	}
	query := r.URL.Query()
	sortParam := query.Get("sort_param")
	if sortParam != "" && sortParam != "rating" && sortParam != "birthday" && sortParam != "name" {
		zapLogger.Errorf("vad sorting param passed: %s", sortParam)
		errText := fmt.Sprintf(`{"error": %s can not be sorting param"}`, sortParam)
		err = response.WriteResponse(w, []byte(errText), http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}
	films, err := h.filmUseCase.GetFilms(sortParam)
	if err != nil {
		zapLogger.Errorf("error in getting films: %s", err)
		errText := `{"error": "internal server error}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			return
		}
		return
	}
	filmsJSON, err := json.Marshal(films)
	if err != nil {
		zapLogger.Errorf("error in marshalling films in json: %s", err)
		errText := `{"error": "internal server error}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			return
		}
		return
	}
	err = response.WriteResponse(w, filmsJSON, http.StatusOK)
	if err != nil {
		zapLogger.Errorf("can not write response: %s", err)
	}
}

// GetFilmByID @Summary Получить фильм по ID
// @Description Получить информацию о фильме по его идентификатору
// @Tags films
// @Accept json
// @Produce json
// @Param FILM_ID path string true "ID актера"
// @Success 200 {object} entity.Film
// @Failure 400 {object} string "Идентификатор фильма передан в неверном формате"
// @Failure 404 {object} string "Фильм не найден"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/film/{FILM_ID} [get]
func (h *FilmHandler) GetFilmByID(w http.ResponseWriter, r *http.Request) {
	zapLogger, err := logger.GetLoggerFromContext(r.Context())
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		err = response.WriteResponse(w, []byte(`"error":"internal error"`), http.StatusInternalServerError)
		if err != nil {
			log.Printf("can not write response: %s", err)
		}
		return
	}
	vars := mux.Vars(r)
	filmID := vars["FILM_ID"]
	filmIDInt, err := strconv.ParseUint(filmID, 10, 64)
	if err != nil {
		zapLogger.Errorf("error in filmID conversion: %s", err)
		errText := fmt.Sprintf(`{"error": "bad format of film id: %s"}`, err)
		err = response.WriteResponse(w, []byte(errText), http.StatusBadRequest)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	film, err := h.filmUseCase.GetFilmByID(filmIDInt)
	if errors.Is(err, usecase.ErrFilmNotFound) {
		zapLogger.Errorf("film with id %d is not found", filmIDInt)
		errText := fmt.Sprintf(`{"error": "film with ID %d is not found"}`, filmIDInt)
		err = response.WriteResponse(w, []byte(errText), http.StatusNotFound)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	if err != nil {
		zapLogger.Errorf("error in getting film: %s", err)
		errText := `{"error": "internal server error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	filmJSON, err := json.Marshal(film)
	if err != nil {
		zapLogger.Errorf("error marshalling response: %s", err)
		errText := `{"error": "internal server error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	err = response.WriteResponse(w, filmJSON, http.StatusOK)
	if err != nil {
		zapLogger.Errorf("error in writing response: %s", err)
	}
}

// AddFilm @Summary Добавление нового фильма
// @Description Данный метод позволяет добавить новый фильм в систему.
// @Tags films
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param body body entity.Film true "Данные о новом фильме"
// @Success 200 {object} entity.Film "Данные добавленного фильма"
// @Failure 400 {object} string "Ошибка в запросе"
// @Failure 401 {object} string "Пользователь не аутентифицирован"
// @Failure 403 {object} string "Запрещено для данного пользователя"
// @Failure 422 {object} string "Ошибка валидации данных"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/admin/film [post]
func (h *FilmHandler) AddFilm(w http.ResponseWriter, r *http.Request) {
	zapLogger, err := logger.GetLoggerFromContext(r.Context())
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		err = response.WriteResponse(w, []byte(`"error":"internal error"`), http.StatusInternalServerError)
		if err != nil {
			log.Printf("can not write response: %s", err)
		}
		return
	}
	filmDTO := &dto.FilmAdd{}
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		zapLogger.Errorf("error in reading request body: %s", err)
		errText := fmt.Sprintf(`{"error": "error in reading request body: %s"}`, err)
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	err = json.Unmarshal(rBody, filmDTO)
	if err != nil {
		zapLogger.Errorf("error in unmarshalling film: %s", err)
		errText := fmt.Sprintf(`{"error": "error in decoding film: %s"}`, err)
		err = response.WriteResponse(w, []byte(errText), http.StatusBadRequest)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}

	if validationErrors := filmDTO.Validate(); len(validationErrors) != 0 {
		var errorsJSON []byte
		errorsJSON, err = json.Marshal(validationErrors)
		if err != nil {
			zapLogger.Errorf("error in marshalling validation errors: %s", err)
			errText := `{"error": "internal server error"}`
			err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
			if err != nil {
				zapLogger.Errorf("error in writing response: %s", err)
			}
			return
		}
		err = response.WriteResponse(w, errorsJSON, http.StatusUnprocessableEntity)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}

	film, actorIDs := filmDTO.GetFilmAndActorIDs()
	addedFilm, err := h.filmUseCase.AddFilm(film, actorIDs)
	if errors.Is(err, usecase.ErrBadFilmAddData) {
		errText := `{"error": "bad add data"}`
		zapLogger.Errorf("error in adding film: %s", err)
		err = response.WriteResponse(w, []byte(errText), http.StatusBadRequest)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	if err != nil {
		errText := `{"error": "internal server error"}`
		zapLogger.Errorf("error in adding film: %s", err)
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	filmJSON, err := json.Marshal(addedFilm)
	if err != nil {
		zapLogger.Errorf("error in marshalling film: %s", err)
		errText := `{"error": "internal error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	err = response.WriteResponse(w, filmJSON, http.StatusOK)
	if err != nil {
		zapLogger.Errorf("error in writing response: %s", err)
	}
}

// UpdateFilm @Summary Обновление информации о фильме
// @Description Данный метод позволяет обновить информацию о фильме.
// @Tags films
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param body body entity.Film true "Данные для обновления фильма"
// @Success 200 {object} entity.Film "Обновленные данные о фильме"
// @Failure 400 {object} string "Ошибка в запросе"
// @Failure 401 {object} string "Пользователь не аутентифицирован"
// @Failure 403 {object} string "Запрещено для данного пользователя"
// @Failure 422 {object} string "Ошибка валидации данных"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/admin/film [put]
func (h *FilmHandler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	zapLogger, err := logger.GetLoggerFromContext(r.Context())
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		err = response.WriteResponse(w, []byte(`"error":"internal error"`), http.StatusInternalServerError)
		if err != nil {
			log.Printf("can not write response: %s", err)
		}
		return
	}
	filmDTO := &dto.FilmUpdate{}
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		zapLogger.Errorf("error in reading request body: %s", err)
		errText := fmt.Sprintf(`{"error": "error in reading request body: %s"}`, err)
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	err = json.Unmarshal(rBody, filmDTO)
	if err != nil {
		zapLogger.Errorf("error in unmarshalling film: %s", err)
		errText := fmt.Sprintf(`{"error": "error in decoding film: %s"}`, err)
		err = response.WriteResponse(w, []byte(errText), http.StatusBadRequest)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}

	if validationErrors := filmDTO.Validate(); len(validationErrors) != 0 {
		var errorsJSON []byte
		errorsJSON, err = json.Marshal(validationErrors)
		if err != nil {
			zapLogger.Errorf("error in marshalling validation errors: %s", err)
			errText := `{"error": "internal server error"}`
			err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
			if err != nil {
				zapLogger.Errorf("error in writing response: %s", err)
			}
			return
		}
		err = response.WriteResponse(w, errorsJSON, http.StatusUnprocessableEntity)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}

	film, actorIDs := filmDTO.GetFilmAndActorIDs()
	err = h.filmUseCase.UpdateFilm(film, actorIDs)
	if errors.Is(err, usecase.ErrBadFilmUpdateData) {
		errText := `{"error": "bad update data"}`
		zapLogger.Errorf("error in updating film: %s", err)
		err = response.WriteResponse(w, []byte(errText), http.StatusBadRequest)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	if err != nil {
		errText := `{"error": "internal server error"}`
		zapLogger.Errorf("error in updating film: %s", err)
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	filmJSON, err := json.Marshal(film)
	if err != nil {
		zapLogger.Errorf("error in marshalling film: %s", err)
		errText := `{"error": "internal error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	err = response.WriteResponse(w, filmJSON, http.StatusOK)
	if err != nil {
		zapLogger.Errorf("error in writing response: %s", err)
	}
}

// GetFilmsBySearch @Summary Получение списка фильмов по поиску
// @Description Данный метод позволяет получить список фильмов, соответствующих поисковому запросу.
// @Tags films
// @Accept json
// @Produce json
// @Param SEARCH_STR path string true "Строка поиска"
// @Success 200 {array} entity.Film "Список фильмов"
// @Failure 404 {object} string "Фильмы не найдены"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/films/search/{SEARCH_STR} [get]
func (h *FilmHandler) GetFilmsBySearch(w http.ResponseWriter, r *http.Request) {
	zapLogger, err := logger.GetLoggerFromContext(r.Context())
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		err = response.WriteResponse(w, []byte(`"error":"internal error"`), http.StatusInternalServerError)
		if err != nil {
			log.Printf("can not write response: %s", err)
		}
		return
	}
	vars := mux.Vars(r)
	searchStr := vars["SEARCH_STR"]
	films, err := h.filmUseCase.GetFilmsBySearch(searchStr)
	if errors.Is(err, usecase.ErrFilmsNotFound) {
		zapLogger.Errorf("no films as a rusult of search %s", searchStr)
		errText := fmt.Sprintf(`{"error": "no films found for search %s"}`, searchStr)
		err = response.WriteResponse(w, []byte(errText), http.StatusNotFound)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	if err != nil {
		errText := `{"error": "internal server error"}`
		zapLogger.Errorf("error in searching films: %s", err)
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	filmsJSON, err := json.Marshal(films)
	if err != nil {
		zapLogger.Errorf("error in marshalling films: %s", err)
		errText := `{"error": "internal error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	err = response.WriteResponse(w, filmsJSON, http.StatusOK)
	if err != nil {
		zapLogger.Errorf("error in writing response: %s", err)
	}
}

// DeleteFilm @Summary Удаление фильма
// @Description Данный метод позволяет удалить фильм по его идентификатору.
// @Tags films
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param FILM_ID path int true "Идентификатор фильма"
// @Success 200 {object} string "Успешное удаление"
// @Failure 400 {object} string "Ошибка в запросе"
// @Failure 401 {object} string "Пользователь не аутентифицирован"
// @Failure 403 {object} string "Запрещено для данного пользователя"
// @Failure 404 {object} string "Фильм не найден"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/admin/film/{FILM_ID} [delete]
func (h *FilmHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	zapLogger, err := logger.GetLoggerFromContext(r.Context())
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		err = response.WriteResponse(w, []byte(`"error":"internal error"`), http.StatusInternalServerError)
		if err != nil {
			log.Printf("can not write response: %s", err)
		}
		return
	}
	vars := mux.Vars(r)
	filmID := vars["FILM_ID"]
	filmIDInt, err := strconv.ParseUint(filmID, 10, 64)
	if err != nil {
		zapLogger.Errorf("error in filmID conversion: %s", err)
		errText := fmt.Sprintf(`{"error": "bad format of film id: %s"}`, err)
		err = response.WriteResponse(w, []byte(errText), http.StatusBadRequest)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	err = h.filmUseCase.DeleteFilm(filmIDInt)
	if errors.Is(err, usecase.ErrFilmNotFound) {
		zapLogger.Errorf("film with id %d is not found", filmIDInt)
		errText := fmt.Sprintf(`{"error": "film with ID %d is not found"}`, filmIDInt)
		err = response.WriteResponse(w, []byte(errText), http.StatusNotFound)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	if err != nil {
		errText := `{"error": "internal server error"}`
		zapLogger.Errorf("error in deleting film: %s", err)
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	result := `{"result": "success"}`
	err = response.WriteResponse(w, []byte(result), http.StatusOK)
	if err != nil {
		zapLogger.Errorf("error in writing response: %s", err)
	}
}
