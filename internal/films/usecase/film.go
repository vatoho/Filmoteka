package usecase

import (
	"github.com/ilyushkaaa/Filmoteka/internal/films/entity"
	"github.com/ilyushkaaa/Filmoteka/internal/films/repo"
)

//go:generate mockgen -source=film.go -destination=film_mock.go -package=usecase FilmUseCase
type FilmUseCase interface {
	GetFilms(sortParam string) ([]entity.Film, error)
	GetFilmByID(filmID uint64) (*entity.Film, error)
	AddFilm(film entity.Film, actorIDs []uint64) (*entity.Film, error)
	UpdateFilm(film entity.Film, actorIDs []uint64) error
	GetFilmsBySearch(searchStr string) ([]entity.Film, error)
	DeleteFilm(ID uint64) error
}

type FilmUseCaseApp struct {
	filmRepo repo.FilmRepo
}

func NewFilmUseCase(filmRepo repo.FilmRepo) *FilmUseCaseApp {
	return &FilmUseCaseApp{
		filmRepo: filmRepo,
	}
}

func (r *FilmUseCaseApp) GetFilms(sortParam string) ([]entity.Film, error) {
	films, err := r.filmRepo.GetFilms(sortParam)
	if err != nil {
		return nil, err
	}
	return films, nil
}

func (r *FilmUseCaseApp) GetFilmByID(filmID uint64) (*entity.Film, error) {
	film, err := r.filmRepo.GetFilmByID(filmID)
	if err != nil {
		return nil, err
	}
	if film == nil {
		return nil, ErrFilmNotFound
	}
	return film, nil
}

func (r *FilmUseCaseApp) AddFilm(film entity.Film, actorIDs []uint64) (*entity.Film, error) {
	filmID, err := r.filmRepo.AddFilm(film, actorIDs)
	if err != nil {
		return nil, err
	}
	if filmID == 0 {
		return nil, ErrBadFilmAddData
	}
	film.ID = filmID
	return &film, nil
}

func (r *FilmUseCaseApp) UpdateFilm(film entity.Film, actorIDs []uint64) error {
	wasUpdated, err := r.filmRepo.UpdateFilm(film, actorIDs)
	if err != nil {
		return err
	}
	if !wasUpdated {
		return ErrBadFilmUpdateData
	}
	return nil
}

func (r *FilmUseCaseApp) GetFilmsBySearch(searchStr string) ([]entity.Film, error) {
	films, err := r.filmRepo.GetFilmsBySearch(searchStr)
	if err != nil {
		return nil, err
	}
	if films == nil || len(films) == 0 {
		return nil, ErrFilmsNotFound
	}
	return films, nil
}

func (r *FilmUseCaseApp) DeleteFilm(ID uint64) error {
	wasDeleted, err := r.filmRepo.DeleteFilm(ID)
	if err != nil {
		return err
	}
	if !wasDeleted {
		return ErrFilmNotFound
	}
	return nil
}
