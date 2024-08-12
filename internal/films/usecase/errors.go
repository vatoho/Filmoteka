package usecase

import "errors"

var (
	ErrFilmsNotFound     = errors.New("films for this search were not found")
	ErrFilmNotFound      = errors.New("film with such id does not exist")
	ErrBadFilmUpdateData = errors.New("invalid data to update film")
	ErrBadFilmAddData    = errors.New("invalid data to add film")
)
