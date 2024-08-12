package middleware

import (
	sessionUseCase "github.com/ilyushkaaa/Filmoteka/internal/session/usecase"
	userUsecase "github.com/ilyushkaaa/Filmoteka/internal/users/usecase"
)

type Middleware struct {
	sessionUseCase sessionUseCase.SessionUseCase
	userUseCase    userUsecase.UserUseCase
}

func NewMiddleware(sessionUseCase sessionUseCase.SessionUseCase, userUseCase userUsecase.UserUseCase) *Middleware {
	return &Middleware{
		sessionUseCase: sessionUseCase,
		userUseCase:    userUseCase,
	}
}
