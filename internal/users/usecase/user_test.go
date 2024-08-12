package usecase

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilyushkaaa/Filmoteka/internal/users/entity"
	"github.com/ilyushkaaa/Filmoteka/internal/users/repo/mock"
	"github.com/ilyushkaaa/Filmoteka/pkg/password_hash"
	"github.com/stretchr/testify/assert"
)

type errorHasher struct{}

func (h *errorHasher) GetHashPassword(password string) (string, error) {
	return "", fmt.Errorf("error")
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testRepo := mock.NewMockUserRepo(ctrl)
	testUseCase := NewUserUseCase(testRepo, &password_hash.SHA256Hasher{})

	var userExpected *entity.User

	errorUseCase := NewUserUseCase(testRepo, &errorHasher{})
	user, err := errorUseCase.Login("aaa", "11111111")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, userExpected, user)

	testRepo.EXPECT().Login("aaa", "ee79976c9380d5e337fc1c095ece8c8f22f91f306ceeb161fa51fecede2c4ba1").
		Return(nil, fmt.Errorf("error"))
	user, err = testUseCase.Login("aaa", "11111111")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, userExpected, user)

	testRepo.EXPECT().Login("aaa", "ee79976c9380d5e337fc1c095ece8c8f22f91f306ceeb161fa51fecede2c4ba1").
		Return(nil, nil)
	user, err = testUseCase.Login("aaa", "11111111")
	assert.Equal(t, ErrBadCredentials, err)
	assert.Equal(t, userExpected, user)

	returnedUser := &entity.User{}
	testRepo.EXPECT().Login("aaa", "ee79976c9380d5e337fc1c095ece8c8f22f91f306ceeb161fa51fecede2c4ba1").
		Return(returnedUser, nil)
	user, err = testUseCase.Login("aaa", "11111111")
	assert.Equal(t, nil, err)
	assert.Equal(t, returnedUser, user)

}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testRepo := mock.NewMockUserRepo(ctrl)
	testUseCase := NewUserUseCase(testRepo, &password_hash.SHA256Hasher{})

	testRepo.EXPECT().GetUserByUsername("aaa").Return(nil, fmt.Errorf("error"))
	user, err := testUseCase.Register("aaa", "11111111")
	var userExpected *entity.User
	assert.NotEqual(t, nil, err)
	assert.Equal(t, userExpected, user)

	returnedUser := &entity.User{}
	testRepo.EXPECT().GetUserByUsername("aaa").Return(returnedUser, nil)
	user, err = testUseCase.Register("aaa", "11111111")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, userExpected, user)

	testRepo.EXPECT().GetUserByUsername("aaa").Return(nil, nil)
	testRepo.EXPECT().Register("aaa",
		"ee79976c9380d5e337fc1c095ece8c8f22f91f306ceeb161fa51fecede2c4ba1").
		Return(nil, fmt.Errorf("error"))
	user, err = testUseCase.Register("aaa", "11111111")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, userExpected, user)

	testRepo.EXPECT().GetUserByUsername("aaa").Return(nil, nil)
	testRepo.EXPECT().Register("aaa",
		"ee79976c9380d5e337fc1c095ece8c8f22f91f306ceeb161fa51fecede2c4ba1").
		Return(returnedUser, nil)
	user, err = testUseCase.Register("aaa", "11111111")
	assert.Equal(t, nil, err)
	assert.Equal(t, returnedUser, user)

	errorUseCase := NewUserUseCase(testRepo, &errorHasher{})
	testRepo.EXPECT().GetUserByUsername("aaa").Return(nil, nil)
	user, err = errorUseCase.Register("aaa", "11111111")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, userExpected, user)

}
