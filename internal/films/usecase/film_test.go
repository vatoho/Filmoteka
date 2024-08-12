package usecase

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilyushkaaa/Filmoteka/internal/films/entity"
	"github.com/ilyushkaaa/Filmoteka/internal/films/repo/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetFilms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testRepo := mock.NewMockFilmRepo(ctrl)
	testUseCase := NewFilmUseCase(testRepo)

	var filmsExpected []entity.Film
	testRepo.EXPECT().GetFilms("").
		Return(nil, fmt.Errorf("error"))
	films, err := testUseCase.GetFilms("")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, filmsExpected, films)

	filmsResult := make([]entity.Film, 0)
	testRepo.EXPECT().GetFilms("").
		Return(filmsResult, nil)
	films, err = testUseCase.GetFilms("")
	assert.Equal(t, nil, err)
	assert.Equal(t, filmsResult, films)

}

func TestGetFilmByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testRepo := mock.NewMockFilmRepo(ctrl)
	testUseCase := NewFilmUseCase(testRepo)

	var id uint64 = 1
	var filmExpected *entity.Film
	testRepo.EXPECT().GetFilmByID(id).
		Return(nil, fmt.Errorf("error"))
	film, err := testUseCase.GetFilmByID(id)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, filmExpected, film)

	filmResult := &entity.Film{
		ID:          id,
		Name:        "aaa",
		Description: "aaa",
	}
	testRepo.EXPECT().GetFilmByID(id).
		Return(filmResult, nil)
	film, err = testUseCase.GetFilmByID(id)
	assert.Equal(t, nil, err)
	assert.Equal(t, filmResult, film)

	testRepo.EXPECT().GetFilmByID(id).
		Return(nil, nil)
	film, err = testUseCase.GetFilmByID(id)
	assert.Equal(t, ErrFilmNotFound, err)
	assert.Equal(t, filmExpected, film)

}

func TestAddFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testRepo := mock.NewMockFilmRepo(ctrl)
	testUseCase := NewFilmUseCase(testRepo)

	var id uint64 = 1
	var idNull uint64 = 0
	var filmExpected *entity.Film
	filmToAdd := entity.Film{}
	actorIDsToAdd := make([]uint64, 0)

	testRepo.EXPECT().AddFilm(filmToAdd, actorIDsToAdd).
		Return(idNull, fmt.Errorf("error"))
	film, err := testUseCase.AddFilm(filmToAdd, actorIDsToAdd)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, filmExpected, film)

	testRepo.EXPECT().AddFilm(filmToAdd, actorIDsToAdd).
		Return(idNull, nil)
	film, err = testUseCase.AddFilm(filmToAdd, actorIDsToAdd)
	assert.Equal(t, ErrBadFilmAddData, err)
	assert.Equal(t, filmExpected, film)

	testRepo.EXPECT().AddFilm(filmToAdd, actorIDsToAdd).
		Return(id, nil)
	film, err = testUseCase.AddFilm(filmToAdd, actorIDsToAdd)
	assert.Equal(t, nil, err)
	filmToAdd.ID = 1
	assert.Equal(t, &filmToAdd, film)

}

func TestUpdateFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testRepo := mock.NewMockFilmRepo(ctrl)
	testUseCase := NewFilmUseCase(testRepo)

	filmToUpdate := entity.Film{}
	actorIDsToAdd := make([]uint64, 0)

	testRepo.EXPECT().UpdateFilm(filmToUpdate, actorIDsToAdd).
		Return(false, fmt.Errorf("error"))
	err := testUseCase.UpdateFilm(filmToUpdate, actorIDsToAdd)
	assert.NotEqual(t, nil, err)

	testRepo.EXPECT().UpdateFilm(filmToUpdate, actorIDsToAdd).
		Return(false, nil)
	err = testUseCase.UpdateFilm(filmToUpdate, actorIDsToAdd)
	assert.Equal(t, ErrBadFilmUpdateData, err)

	testRepo.EXPECT().UpdateFilm(filmToUpdate, actorIDsToAdd).
		Return(true, nil)
	err = testUseCase.UpdateFilm(filmToUpdate, actorIDsToAdd)
	assert.Equal(t, nil, err)
}

func TestDeleteFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testRepo := mock.NewMockFilmRepo(ctrl)
	testUseCase := NewFilmUseCase(testRepo)

	var id uint64 = 1
	testRepo.EXPECT().DeleteFilm(id).
		Return(false, fmt.Errorf("error"))
	err := testUseCase.DeleteFilm(id)
	assert.NotEqual(t, nil, err)

	testRepo.EXPECT().DeleteFilm(id).
		Return(false, nil)
	err = testUseCase.DeleteFilm(id)
	assert.Equal(t, ErrFilmNotFound, err)

	testRepo.EXPECT().DeleteFilm(id).
		Return(true, nil)
	err = testUseCase.DeleteFilm(id)
	assert.Equal(t, nil, err)
}
