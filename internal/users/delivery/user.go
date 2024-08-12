package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ilyushkaaa/Filmoteka/internal/dto"
	"github.com/ilyushkaaa/Filmoteka/internal/middleware"
	usecaseSession "github.com/ilyushkaaa/Filmoteka/internal/session/usecase"
	"github.com/ilyushkaaa/Filmoteka/internal/users/entity"
	usecaseUser "github.com/ilyushkaaa/Filmoteka/internal/users/usecase"
	"github.com/ilyushkaaa/Filmoteka/pkg/logger"
	"github.com/ilyushkaaa/Filmoteka/pkg/response"
	"go.uber.org/zap"
)

type UserHandler struct {
	userUseCase    usecaseUser.UserUseCase
	sessionUseCase usecaseSession.SessionUseCase
}

func NewUserHandler(userUseCase usecaseUser.UserUseCase, sessionUseCase usecaseSession.SessionUseCase) *UserHandler {
	return &UserHandler{
		userUseCase:    userUseCase,
		sessionUseCase: sessionUseCase,
	}
}

// Login @Summary Вход пользователя
// @Description Данный метод позволяет пользователям войти в систему, используя свои учетные данные.
// @Tags users
// @Accept json
// @Produce json
// @Param body body dto.AuthRequest true "Данные пользователя для входа"
// @Success 200 {object} dto.AuthResponse "Успешный вход, получен идентификатор сессии"
// @Failure 401 {object} string "Неверные учетные данные"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/login [post]
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	zapLogger, err := logger.GetLoggerFromContext(ctx)
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		err = response.WriteResponse(w, []byte(`"error":"internal error"`), http.StatusInternalServerError)
		if err != nil {
			log.Printf("can not write response: %s", err)
		}
		return
	}
	userFromLoginForm, err := checkRequestFormat(zapLogger, w, r)
	if err != nil || userFromLoginForm == nil {
		return
	}
	loggedInUser, err := uh.userUseCase.Login(userFromLoginForm.Username, userFromLoginForm.Password)

	if errors.Is(err, usecaseUser.ErrBadCredentials) {
		zapLogger.Errorf("bad credentials were eneterd")
		err = response.WriteResponse(w, []byte(`{"error": "bad username or password"}`), http.StatusUnauthorized)
		if err != nil {
			zapLogger.Errorf("can not write response: %s", err)
		}
		return
	}
	if err != nil {
		zapLogger.Errorf("internal error in logging in user: %s", err)
		errText := fmt.Sprintf(`{"error": "error in getting user by login and password: %s"}`, err)
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("can not write response: %s", err)
		}
		return
	}

	uh.HandleGetSessionID(w, loggedInUser, zapLogger)

}

// Register @Summary Регистрация пользователя
// @Description Данный метод позволяет новым пользователям зарегистрироваться в системе.
// @Tags users
// @Accept json
// @Produce json
// @Param body body dto.AuthRequest true "Данные нового пользователя"
// @Success 200 {object} dto.AuthResponse "Успешная регистрация, получен идентификатор сессии"
// @Failure 422 {object} string "Пользователь уже существует"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/register [post]
func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	zapLogger, err := logger.GetLoggerFromContext(ctx)
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		err = response.WriteResponse(w, []byte(`"error":"internal error"`), http.StatusInternalServerError)
		if err != nil {
			log.Printf("can not write response: %s", err)
		}
		return
	}
	userFromLoginForm, err := checkRequestFormat(zapLogger, w, r)
	if err != nil || userFromLoginForm == nil {
		return
	}
	newUser, err := uh.userUseCase.Register(userFromLoginForm.Username, userFromLoginForm.Password)
	if errors.Is(err, usecaseUser.ErrUserAlreadyExists) {
		zapLogger.Errorf("user with username %s alredy exists", userFromLoginForm.Username)
		err = response.WriteResponse(w, []byte(`{"error": "user already exists"}`), http.StatusUnprocessableEntity)
		if err != nil {
			zapLogger.Errorf("can not write response: %s", err)
		}
		return
	}
	if err != nil {
		zapLogger.Errorf("internal error in register: %s", err)
		errText := `{"error": "internal error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("can not write response: %s", err)
		}
		return
	}
	uh.HandleGetSessionID(w, newUser, zapLogger)
}

func (uh *UserHandler) HandleGetSessionID(w http.ResponseWriter, newUser *entity.User, zapLogger *zap.SugaredLogger) {
	sessionID, err := uh.sessionUseCase.CreateSession(newUser.ID)
	if err != nil {
		zapLogger.Errorf("internal error in getting session id: %s", err)
		errText := `{"error": "internal error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("can not write response: %s", err)
		}
		return
	}
	resp := dto.AuthResponse{
		SessionID: sessionID,
	}
	sessionIDJSON, err := json.Marshal(&resp)
	if err != nil {
		zapLogger.Errorf("error in marshalling session: %s", err)
		errText := `{"error": "internal error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("can not write response: %s", err)
		}
		return
	}
	zapLogger.Infof("new session id: %s", sessionID)
	err = response.WriteResponse(w, sessionIDJSON, http.StatusOK)
	if err != nil {
		zapLogger.Errorf("can not write response: %s", err)
	}
}

func checkRequestFormat(zapLogger *zap.SugaredLogger, w http.ResponseWriter, r *http.Request) (*dto.AuthRequest, error) {
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		zapLogger.Errorf("error in reading request: %s", err)
		errText := fmt.Sprintf(`{"error": "error in reading request body: %s"}`, err)
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)

		if err != nil {
			zapLogger.Errorf("can not write response: %s", err)
		}

		return nil, err
	}
	userFromLoginForm := &dto.AuthRequest{}
	err = json.Unmarshal(rBody, userFromLoginForm)
	if err != nil {
		zapLogger.Errorf("error in unmarshalling user: %s", err)
		errText := fmt.Sprintf(`{"error": "error in decoding user: %s"}`, err)
		err = response.WriteResponse(w, []byte(errText), http.StatusUnauthorized)
		if err != nil {
			zapLogger.Errorf("can not write response: %s", err)
		}
		return nil, err
	}
	if validationErrors := userFromLoginForm.Validate(); len(validationErrors) != 0 {
		errorsJSON, err := json.Marshal(validationErrors)
		if err != nil {
			zapLogger.Errorf("error in marshalling errors: %s", err)
			errText := `{"error": "internal error"}`
			err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
			if err != nil {
				zapLogger.Errorf("can not write response: %s", err)
			}
			return nil, err
		}
		zapLogger.Errorf("login form did not pass validation: %s", err)
		err = response.WriteResponse(w, errorsJSON, http.StatusUnauthorized)
		if err != nil {
			zapLogger.Errorf("can not write response: %s", err)
		}
		return nil, err
	}
	return userFromLoginForm, nil
}

// Logout @Summary Выход пользователя
// @Description Данный метод позволяет пользователям выйти из системы, завершая сеанс.
// @Tags users
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 200 {object} string "Успешный выход"
// @Failure 401 {object} string "Неверный или отсутствующий токен аутентификации"
// @Failure 404 {object} string "Сеанс не найден"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/v1/logout [post]
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	zapLogger, err := logger.GetLoggerFromContext(ctx)
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		err = response.WriteResponse(w, []byte(`"error": "internal error"`), http.StatusInternalServerError)
		if err != nil {
			log.Printf("can not write response: %s", err)
		}
		return
	}
	sessionID, ok := ctx.Value(middleware.MySessionIDKey).(string)
	if !ok {
		zapLogger.Errorf("can not get sesssion id from context: %s", err)
		err = response.WriteResponse(w, []byte(`{"error": "internal error"}`), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("can not write response: %s", err)
		}
		return
	}
	isDeleted, err := uh.sessionUseCase.DeleteSession(sessionID)
	if err != nil {
		zapLogger.Errorf("error in deleting session: %s", err)
		errText := `{"error": "internal error"}`
		err = response.WriteResponse(w, []byte(errText), http.StatusInternalServerError)
		if err != nil {
			zapLogger.Errorf("can not write response: %s", err)
		}
		return
	}
	if !isDeleted {
		zapLogger.Errorf("session with id %s is not found:", sessionID)
		errText := fmt.Sprintf(`{"error": "no session with session id: %s"}`, sessionID)
		err = response.WriteResponse(w, []byte(errText), http.StatusNotFound)
		if err != nil {
			zapLogger.Errorf("can not write response: %s", err)
		}
		return
	}
	message := `{"result":"success"}`
	err = response.WriteResponse(w, []byte(message), http.StatusOK)
	if err != nil {
		zapLogger.Errorf("can not write response: %s", err)
	}
}
