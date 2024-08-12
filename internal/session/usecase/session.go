package usecase

import (
	"github.com/google/uuid"
	"github.com/ilyushkaaa/Filmoteka/internal/session/entity"
	"github.com/ilyushkaaa/Filmoteka/internal/session/repo"
)
//go:generate mockgen -source=session.go -destination=session_mock.go -package=usecase SessionUseCase
type SessionUseCase interface {
	CreateSession(userID uint64) (string, error)
	GetSession(sessionID string) (*entity.Session, error)
	DeleteSession(sessionID string) (bool, error)
}

type SessionUseCaseApp struct {
	sessionRepo repo.SessionRepo
}

func NewSessionUseCase(sessionRepo repo.SessionRepo) *SessionUseCaseApp {
	return &SessionUseCaseApp{
		sessionRepo: sessionRepo,
	}
}

func (su *SessionUseCaseApp) CreateSession(userID uint64) (string, error) {
	newSession := &entity.Session{
		ID:     uuid.New().String(),
		UserID: userID,
	}
	err := su.sessionRepo.CreateSession(newSession)
	return newSession.ID, err
}

func (su *SessionUseCaseApp) GetSession(sessionID string) (*entity.Session, error) {
	session, err := su.sessionRepo.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrNoSession
	}
	return session, nil
}

func (su *SessionUseCaseApp) DeleteSession(sessionID string) (bool, error) {
	isDeleted, err := su.sessionRepo.DeleteSession(sessionID)
	if err != nil {
		return false, err
	}
	if !isDeleted {
		return false, ErrNoSession
	}
	return true, nil
}
