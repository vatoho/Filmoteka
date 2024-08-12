package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ilyushkaaa/Filmoteka/internal/actors/usecase"
	"github.com/ilyushkaaa/Filmoteka/internal/dto"
	"github.com/ilyushkaaa/Filmoteka/pkg/logger"
	"github.com/ilyushkaaa/Filmoteka/pkg/response"
)

type ActorHandler struct {
	actorUseCase usecase.ActorUseCase
}

func NewActorHandler(actorUseCase usecase.ActorUseCase) *ActorHandler {
	return &ActorHandler{
		actorUseCase: actorUseCase,
	}
}

// GetActors @Summary Получить всех актеров
// @Description Получить список всех актеров
// @Tags actors
// @Accept json
// @Produce json
// @Success 200 {array} dto.ActorWithFilms
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/actors [get]
func (h *ActorHandler) GetActors(w http.ResponseWriter, r *http.Request) {
	zapLogger, err := logger.GetLoggerFromContext(r.Context())
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		err = response.WriteResponse(w, []byte(`"error":"internal error"`), http.StatusInternalServerError)
		if err != nil {
			log.Printf("can not write response: %s", err)
		}
		return
	}

	actors, err := h.actorUseCase.GetActors()
	if err != nil {
		zapLogger.Errorf("error in getting actors: %s", err)
		errText := `{"error": "internal server error}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			return
		}
		return
	}
	actorsJSON, err := json.Marshal(actors)
	if err != nil {
		zapLogger.Errorf("error in marshalling actors in json: %s", err)
		errText := `{"error": "internal error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			return
		}
		return
	}
	err = response.WriteResponse(w, actorsJSON, http.StatusOK)
	if err != nil {
		zapLogger.Errorf("can not write response: %s", err)
	}
}

// GetActorByID @Summary Получить актера по ID
// @Description Получить информацию об актере по его ID с информацией о фильмах, в которых он снимался
// @Tags actors
// @Accept json
// @Produce json
// @Param ACTOR_ID path string true "ID актера"
// @Success 200 {object} dto.ActorWithFilms
// @Failure 400 {object} string "Идентификатор актера передан в неверном формате"
// @Failure 404 {object} string "Актёр не найден"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/actor/{ACTOR_ID} [get]
func (h *ActorHandler) GetActorByID(w http.ResponseWriter, r *http.Request) {
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
	actorID := vars["ACTOR_ID"]
	actorIDInt, err := strconv.ParseUint(actorID, 10, 64)
	if err != nil {
		zapLogger.Errorf("error in actor id conversion: %s", err)
		errText := fmt.Sprintf(`{"error": "bad format of actor id: %s"}`, err)
		err = response.WriteResponse(w, []byte(errText), http.StatusBadRequest)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	actor, err := h.actorUseCase.GetActorByID(actorIDInt)
	if errors.Is(err, usecase.ErrActorNotFound) {
		zapLogger.Errorf("actor with id %d is not found", actorIDInt)
		errText := fmt.Sprintf(`{"error": "actor with ID %d is not found"}`, actorIDInt)
		err = response.WriteResponse(w, []byte(errText), http.StatusNotFound)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	if err != nil {
		zapLogger.Errorf("error in getting actor: %s", err)
		errText := `{"error": "internal server error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	actorJSON, err := json.Marshal(actor)
	if err != nil {
		zapLogger.Errorf("error marshalling response: %s", err)
		errText := `{"error": "internal server error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	err = response.WriteResponse(w, actorJSON, http.StatusOK)
	if err != nil {
		zapLogger.Errorf("error in writing response: %s", err)
	}
}

// AddActor @Summary Добавление нового актера
// @Description Данный метод позволяет добавить нового актера в систему.
// @Tags actors
// @Accept json
// @Produce json
// @SecurityRequirement CookieAuth
// @Param body body dto.ActorAdd true "Данные о новом актере"
// @Success 200 {object} dto.ActorAdd "Данные добавленного актера"
// @Failure 400 {object} string "Ошибка в запросе"
// @Failure 401 {object} string "Пользователь не аутентифицирован"
// @Failure 403 {object} string "Запрещено для данного пользователя"
// @Failure 422 {object} string "Ошибка валидации данных"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/admin/actor [post]
func (h *ActorHandler) AddActor(w http.ResponseWriter, r *http.Request) {
	zapLogger, err := logger.GetLoggerFromContext(r.Context())
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		err = response.WriteResponse(w, []byte(`"error":"internal error"`), http.StatusInternalServerError)
		if err != nil {
			log.Printf("can not write response: %s", err)
		}
		return
	}
	actorDTO := &dto.ActorAdd{}
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
	err = json.Unmarshal(rBody, actorDTO)
	if err != nil {
		zapLogger.Errorf("error in unmarshalling actor: %s", err)
		errText := fmt.Sprintf(`{"error": "error in decoding actor: %s"}`, err)
		err = response.WriteResponse(w, []byte(errText), http.StatusBadRequest)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}

	if validationErrors := actorDTO.Validate(); len(validationErrors) != 0 {
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

	actor := actorDTO.Convert()
	addedActor, err := h.actorUseCase.AddActor(actor)
	if err != nil {
		errText := `{"error": "internal server error"}`
		zapLogger.Errorf("error in adding actor: %s", err)
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	actorJSON, err := json.Marshal(addedActor)
	if err != nil {
		zapLogger.Errorf("error in marshalling actor: %s", err)
		errText := `{"error": "internal error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	err = response.WriteResponse(w, actorJSON, http.StatusOK)
	if err != nil {
		zapLogger.Errorf("error in writing response: %s", err)
	}
}

// UpdateActor @Summary Обновление информации об актере
// @Description Данный метод позволяет обновить информацию об актере.
// @Tags actors
// @Accept json
// @Produce json
// @SecurityRequirement CookieAuth
// @Param body body dto.ActorUpdate true "Данные для обновления актера"
// @Success 200 {object} dto.ActorUpdate "Обновленные данные об актере"
// @Failure 400 {object} string "Ошибка в запросе"
// @Failure 401 {object} string "Пользователь не аутентифицирован"
// @Failure 403 {object} string "Запрещено для данного пользователя"
// @Failure 422 {object} string "Ошибка валидации данных"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/admin/actor [put]
func (h *ActorHandler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	zapLogger, err := logger.GetLoggerFromContext(r.Context())
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		err = response.WriteResponse(w, []byte(`"error":"internal error"`), http.StatusInternalServerError)
		if err != nil {
			log.Printf("can not write response: %s", err)
		}
		return
	}
	actorDTO := &dto.ActorUpdate{}
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
	err = json.Unmarshal(rBody, actorDTO)
	if err != nil {
		zapLogger.Errorf("error in unmarshalling actor: %s", err)
		errText := fmt.Sprintf(`{"error": "error in decoding actor: %s"}`, err)
		err = response.WriteResponse(w, []byte(errText), http.StatusBadRequest)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}

	if validationErrors := actorDTO.Validate(); len(validationErrors) != 0 {
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

	actor := actorDTO.Convert()
	err = h.actorUseCase.UpdateActor(actor)
	if errors.Is(err, usecase.ErrActorNotFound) {
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
		zapLogger.Errorf("error in updating actor: %s", err)
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	actorJSON, err := json.Marshal(actor)
	if err != nil {
		zapLogger.Errorf("error in marshalling actor: %s", err)
		errText := `{"error": "internal error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	err = response.WriteResponse(w, actorJSON, http.StatusOK)
	if err != nil {
		zapLogger.Errorf("error in writing response: %s", err)
	}
}

// DeleteActor @Summary Удаление актера
// @Description Данный метод позволяет удалить актера по его идентификатору.
// @Tags actors
// @Accept json
// @Produce json
// @SecurityRequirement CookieAuth
// @Param ACTOR_ID path int true "Идентификатор актера"
// @Success 200 {object} string "Успешное удаление"
// @Failure 400 {object} string "Ошибка в запросе"
// @Failure 401 {object} string "Пользователь не аутентифицирован"
// @Failure 403 {object} string "Запрещено для данного пользователя"
// @Failure 404 {object} string "Актер не найден"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/admin/actor/{ACTOR_ID} [delete]
func (h *ActorHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
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
	actorID := vars["ACTOR_ID"]
	actorIDInt, err := strconv.ParseUint(actorID, 10, 64)
	if err != nil {
		zapLogger.Errorf("error in actor id conversion: %s", err)
		errText := fmt.Sprintf(`{"error": "bad format of actor id: %s"}`, err)
		err = response.WriteResponse(w, []byte(errText), http.StatusBadRequest)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	err = h.actorUseCase.DeleteActor(actorIDInt)
	if errors.Is(err, usecase.ErrActorNotFound) {
		zapLogger.Errorf("actor with id %d is not found", actorIDInt)
		errText := fmt.Sprintf(`{"error": "actor with ID %d is not found"}`, actorIDInt)
		err = response.WriteResponse(w, []byte(errText), http.StatusNotFound)
		if err != nil {
			zapLogger.Errorf("error in writing response: %s", err)
		}
		return
	}
	if err != nil {
		errText := `{"error": "internal server error"}`
		zapLogger.Errorf("error in deleting actor: %s", err)
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
