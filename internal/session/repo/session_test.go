package repo

import (
	"fmt"
	"testing"

	"github.com/ilyushkaaa/Filmoteka/internal/session/entity"
	"github.com/stretchr/testify/assert"
)

type MockRedisConn struct{}

func (m *MockRedisConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	sessionID := args[0].(string)
	if commandName == "GET" || commandName == "EXISTS" || commandName == "DEL" {
		switch sessionID {
		case "valid_session":
			return int64(123), nil
		}
	}
	if commandName == "SET" {
		switch sessionID {
		case "valid_session":
			return "OK", nil
		}
	}
	return nil, fmt.Errorf("err")
}

func TestCreateSession(t *testing.T) {
	red := &MockRedisConn{}
	repo := NewSessionRepo(red)

	sessionInvalid := &entity.Session{
		ID:     "invalid",
		UserID: uint64(1),
	}
	err := repo.CreateSession(sessionInvalid)
	assert.Error(t, err)

	sessionValid := &entity.Session{
		ID:     "valid_session",
		UserID: uint64(1),
	}
	err = repo.CreateSession(sessionValid)
	assert.NoError(t, err)
}

func TestGetSession(t *testing.T) {
	red := &MockRedisConn{}
	repo := NewSessionRepo(red)

	_, err := repo.GetSession("valid_session")
	assert.NoError(t, err)

	_, err = repo.GetSession("invalid_session")
	assert.Error(t, err)
}

func TestDeleteSession(t *testing.T) {
	red := &MockRedisConn{}
	repo := NewSessionRepo(red)

	_, err := repo.DeleteSession("valid_session")
	assert.NoError(t, err)
}

func (m *MockRedisConn) Close() error {
	return nil
}

func (m *MockRedisConn) Err() error {
	return nil
}

func (m *MockRedisConn) Send(commandName string, args ...interface{}) error {
	return nil
}

func (m *MockRedisConn) Flush() error {
	return nil
}

func (m *MockRedisConn) Receive() (reply interface{}, err error) {
	return nil, nil
}
