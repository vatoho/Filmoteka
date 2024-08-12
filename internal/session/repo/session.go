package repo

import (
	"github.com/gomodule/redigo/redis"
	"github.com/ilyushkaaa/Filmoteka/internal/session/entity"
)

//go:generate mockgen -source=session.go -destination=session_mock.go -package=repo SessionRepo
type SessionRepo interface {
	CreateSession(session *entity.Session) error
	GetSession(sessionID string) (*entity.Session, error)
	DeleteSession(sessionID string) (bool, error)
}

type SessionRepoRedis struct {
	redisConn  redis.Conn
	expireTime int
}

func NewSessionRepo(redisConn redis.Conn) *SessionRepoRedis {
	return &SessionRepoRedis{
		redisConn:  redisConn,
		expireTime: 24 * 60 * 60,
	}
}

func (s *SessionRepoRedis) CreateSession(session *entity.Session) error {
	result, err := redis.String(s.redisConn.Do("SET", session.ID, session.UserID, "EX", s.expireTime))
	if err != nil || result != "OK" {
		return err
	}
	return nil
}

func (s *SessionRepoRedis) GetSession(sessionID string) (*entity.Session, error) {
	userID, err := redis.Uint64(s.redisConn.Do("GET", sessionID))
	if err != nil {
		return nil, err
	}
	session := &entity.Session{
		ID:     sessionID,
		UserID: userID,
	}
	return session, nil

}

func (s *SessionRepoRedis) DeleteSession(sessionID string) (bool, error) {
	exists, err := redis.Bool(s.redisConn.Do("EXISTS", sessionID))
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	_, err = s.redisConn.Do("DEL", sessionID)
	if err != nil {
		return false, err
	}
	return true, nil
}
