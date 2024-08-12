package usecase

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilyushkaaa/Filmoteka/internal/actors/entity"
	"github.com/ilyushkaaa/Filmoteka/internal/actors/repo/mock"
	"github.com/ilyushkaaa/Filmoteka/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestGetActors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testRepo := mock.NewMockActorRepo(ctrl)
	testUseCase := NewActorUseCase(testRepo)

	var actorWithFilmsExpected []dto.ActorWithFilms
	testRepo.EXPECT().GetActors().
		Return(nil, fmt.Errorf("error"))
	actors, err := testUseCase.GetActors()
	assert.NotEqual(t, nil, err)
	assert.Equal(t, actorWithFilmsExpected, actors)

	actorWithFilmsResult := make([]dto.ActorWithFilms, 0)
	testRepo.EXPECT().GetActors().
		Return(actorWithFilmsResult, nil)
	actors, err = testUseCase.GetActors()
	assert.Equal(t, nil, err)
	assert.Equal(t, actorWithFilmsResult, actors)

}

func TestGetActorByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testRepo := mock.NewMockActorRepo(ctrl)
	testUseCase := NewActorUseCase(testRepo)

	var id uint64 = 1
	var actorWithFilmsExpected *dto.ActorWithFilms
	testRepo.EXPECT().GetActorByID(id).
		Return(nil, fmt.Errorf("error"))
	actors, err := testUseCase.GetActorByID(id)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, actorWithFilmsExpected, actors)

	testRepo.EXPECT().GetActorByID(id).
		Return(nil, nil)
	actors, err = testUseCase.GetActorByID(id)
	assert.Equal(t, ErrActorNotFound, err)
	assert.Equal(t, actorWithFilmsExpected, actors)

	actorWithFilmsResult := &dto.ActorWithFilms{}
	testRepo.EXPECT().GetActorByID(id).
		Return(actorWithFilmsResult, nil)
	actors, err = testUseCase.GetActorByID(id)
	assert.Equal(t, nil, err)
	assert.Equal(t, actorWithFilmsResult, actors)
}

func TestAddActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testRepo := mock.NewMockActorRepo(ctrl)
	testUseCase := NewActorUseCase(testRepo)

	var id uint64 = 1
	actorToAdd := entity.Actor{}
	var actorExpected *entity.Actor
	testRepo.EXPECT().AddActor(actorToAdd).
		Return(id, fmt.Errorf("error"))
	actor, err := testUseCase.AddActor(actorToAdd)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, actorExpected, actor)

	testRepo.EXPECT().AddActor(actorToAdd).
		Return(id, nil)
	actor, err = testUseCase.AddActor(actorToAdd)
	actorToAdd.ID = id
	assert.Equal(t, nil, err)
	assert.Equal(t, &actorToAdd, actor)
}

func TestUpdateActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testRepo := mock.NewMockActorRepo(ctrl)
	testUseCase := NewActorUseCase(testRepo)

	actorToUpdate := entity.Actor{}
	testRepo.EXPECT().UpdateActor(actorToUpdate).
		Return(false, fmt.Errorf("error"))
	err := testUseCase.UpdateActor(actorToUpdate)
	assert.NotEqual(t, nil, err)

	testRepo.EXPECT().UpdateActor(actorToUpdate).
		Return(false, nil)
	err = testUseCase.UpdateActor(actorToUpdate)
	assert.Equal(t, ErrActorNotFound, err)
}
