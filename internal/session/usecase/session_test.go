package usecase

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilyushkaaa/Filmoteka/internal/session/entity"
	"github.com/ilyushkaaa/Filmoteka/internal/session/repo/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testRepo := mock.NewMockSessionRepo(ctrl)
	testUseCase := NewSessionUseCase(testRepo)

	var sessionExpected *entity.Session
	testRepo.EXPECT().GetSession("qqqq").
		Return(nil, fmt.Errorf("error"))
	session, err := testUseCase.GetSession("qqqq")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, sessionExpected, session)

	testRepo.EXPECT().GetSession("qqqq").
		Return(nil, nil)
	session, err = testUseCase.GetSession("qqqq")
	assert.Equal(t, ErrNoSession, err)
	assert.Equal(t, sessionExpected, session)

	sessionResult := &entity.Session{}
	testRepo.EXPECT().GetSession("qqqq").
		Return(sessionResult, nil)
	session, err = testUseCase.GetSession("qqqq")
	assert.Equal(t, nil, err)
	assert.Equal(t, sessionResult, session)

}

func TestDeleteSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testRepo := mock.NewMockSessionRepo(ctrl)
	testUseCase := NewSessionUseCase(testRepo)

	testRepo.EXPECT().DeleteSession("qqqq").
		Return(false, fmt.Errorf("error"))
	wasDeleted, err := testUseCase.DeleteSession("qqqq")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, false, wasDeleted)

	testRepo.EXPECT().DeleteSession("qqqq").
		Return(false, nil)
	wasDeleted, err = testUseCase.DeleteSession("qqqq")
	assert.Equal(t, ErrNoSession, err)
	assert.Equal(t, false, wasDeleted)

	testRepo.EXPECT().DeleteSession("qqqq").
		Return(true, nil)
	wasDeleted, err = testUseCase.DeleteSession("qqqq")
	assert.Equal(t, nil, err)
	assert.Equal(t, true, wasDeleted)
}
