package repo

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/ilyushkaaa/Filmoteka/internal/films/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetFilmByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	repo := &FilmRepoPG{db: db, zapLogger: nil}

	expectedFilm := &entity.Film{
		ID:            1,
		Name:          "Film 1",
		Description:   "Description 1",
		DateOfRelease: time.Time{}.Add(time.Hour),
		Rating:        8.5,
	}

	mock.ExpectQuery("SELECT id, name, description, date_of_release, rating FROM films WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "date_of_release", "rating"}).
			AddRow(1, "Film 1", "Description 1", time.Time{}.Add(time.Hour), 8.5))

	film, err := repo.GetFilmByID(1)

	assert.NoError(t, err, "unexpected error")
	assert.Equal(t, expectedFilm, film, "films do not match expected")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT id, name, description, date_of_release, rating FROM films WHERE id = ?").
		WithArgs(1).
		WillReturnError(fmt.Errorf("error"))

	film, err = repo.GetFilmByID(1)

	var nilFilm *entity.Film
	assert.Error(t, err)
	assert.Equal(t, nilFilm, film)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetFilms(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	repo := &FilmRepoPG{db: db, zapLogger: nil}

	expectedFilms := []entity.Film{
		{ID: 1, Name: "Film 1", Description: "Description 1", DateOfRelease: time.Time{}.Add(time.Hour), Rating: 8.5},
		{ID: 2, Name: "Film 2", Description: "Description 2", DateOfRelease: time.Time{}.Add(time.Hour), Rating: 7.8},
	}

	mock.ExpectQuery("SELECT f.id, f.name, f.description, f.date_of_release, f.rating FROM films f ORDER BY (.+) DESC").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "date_of_release", "rating"}).
			AddRow(1, "Film 1", "Description 1", time.Time{}.Add(time.Hour), 8.5).
			AddRow(2, "Film 2", "Description 2", time.Time{}.Add(time.Hour), 7.8))

	films, err := repo.GetFilms("rating")
	assert.NoError(t, err)
	assert.Equal(t, expectedFilms, films)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT f.id, f.name, f.description, f.date_of_release, f.rating FROM films f ORDER BY (.+) DESC").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "date_of_release", "rating"}).
			AddRow(1, "Film 1", "Description 1", time.Time{}.Add(time.Hour), 8.5).
			AddRow(2, "Film 2", "Description 2", time.Time{}.Add(time.Hour), 7.8))

	films, err = repo.GetFilms("")
	assert.NoError(t, err)
	assert.Equal(t, expectedFilms, films)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT f.id, f.name, f.description, f.date_of_release, f.rating FROM films f ORDER BY (.+) DESC").
		WillReturnError(fmt.Errorf("error"))

	var nilFilms []entity.Film
	films, err = repo.GetFilms("")
	assert.Error(t, err)
	assert.Equal(t, nilFilms, films)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT f.id, f.name, f.description, f.date_of_release, f.rating FROM films f ORDER BY (.+) DESC").
		WillReturnError(sql.ErrNoRows)

	films, err = repo.GetFilms("")
	assert.NoError(t, err)
	assert.Equal(t, nilFilms, films)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAddFilm(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	repo := &FilmRepoPG{db: db, zapLogger: zap.NewNop().Sugar()}

	expectedLastInsertID := uint64(1)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO films").
		WithArgs("Film 1", "Description 1", time.Time{}.Add(time.Hour), 8.5).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedLastInsertID))

	for _, actorID := range []uint64{1, 2} {
		mock.ExpectQuery("SELECT id FROM actors WHERE id = ?").
			WithArgs(actorID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(actorID))
		mock.ExpectExec("INSERT INTO film_actors").
			WithArgs(expectedLastInsertID, actorID).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	mock.ExpectCommit()

	lastInsertID, err := repo.AddFilm(entity.Film{Name: "Film 1", Description: "Description 1", DateOfRelease: time.Time{}.Add(time.Hour), Rating: 8.5}, []uint64{1, 2})

	assert.NoError(t, err, "unexpected error")
	assert.Equal(t, expectedLastInsertID, lastInsertID, "last insert ID does not match expected")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO films").
		WithArgs("Film 1", "Description 1", time.Time{}.Add(time.Hour), 8.5).
		WillReturnError(fmt.Errorf("error"))

	mock.ExpectRollback()

	lastInsertID, err = repo.AddFilm(entity.Film{Name: "Film 1", Description: "Description 1", DateOfRelease: time.Time{}.Add(time.Hour), Rating: 8.5}, []uint64{1, 2})

	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO films").
		WithArgs("Film 1", "Description 1", time.Time{}.Add(time.Hour), 8.5).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedLastInsertID))

	var actID uint64 = 1

	mock.ExpectQuery("SELECT id FROM actors WHERE id = ?").
		WithArgs(actID).
		WillReturnError(fmt.Errorf("err"))

	mock.ExpectRollback()

	lastInsertID, err = repo.AddFilm(entity.Film{Name: "Film 1", Description: "Description 1", DateOfRelease: time.Time{}.Add(time.Hour), Rating: 8.5}, []uint64{1, 2})

	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateFilm(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	repo := &FilmRepoPG{db: db, zapLogger: nil}

	filmID := uint64(1)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE films").
		WithArgs("Updated Film", "Updated Description", time.Time{}.Add(time.Hour), 9.0, filmID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM film_actors").
		WithArgs(filmID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	for _, actorID := range []uint64{1, 2} {
		mock.ExpectQuery("SELECT id FROM actors WHERE id = ?").
			WithArgs(actorID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(actorID))
		mock.ExpectExec("INSERT INTO film_actors").
			WithArgs(filmID, actorID).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	mock.ExpectCommit()

	updated, err := repo.UpdateFilm(entity.Film{ID: filmID, Name: "Updated Film", Description: "Updated Description", DateOfRelease: time.Time{}.Add(time.Hour), Rating: 9.0}, []uint64{1, 2})

	assert.NoError(t, err, "unexpected error")
	assert.True(t, updated, "film was not updated successfully")
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE films").
		WithArgs("Updated Film", "Updated Description", time.Time{}.Add(time.Hour), 9.0, filmID).
		WillReturnError(fmt.Errorf("error"))

	mock.ExpectRollback()

	updated, err = repo.UpdateFilm(entity.Film{ID: filmID, Name: "Updated Film", Description: "Updated Description", DateOfRelease: time.Time{}.Add(time.Hour), Rating: 9.0}, []uint64{1, 2})

	assert.Error(t, err)
	assert.False(t, updated)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetFilmsBySearch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	repo := &FilmRepoPG{db: db, zapLogger: zap.NewNop().Sugar()}

	searchStr := "Film"

	mock.ExpectQuery("SELECT DISTINCT").
		WithArgs(searchStr).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "date_of_release", "rating"}).
			AddRow(1, "Film 1", "Description 1", time.Time{}.Add(time.Hour), 8.0).
			AddRow(2, "Film 2", "Description 2", time.Time{}.Add(time.Hour), 7.5))

	films, err := repo.GetFilmsBySearch(searchStr)

	assert.NoError(t, err, "unexpected error")
	assert.NotNil(t, films, "films list is nil")
	assert.Equal(t, 2, len(films), "unexpected number of films returned")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT DISTINCT").
		WithArgs(searchStr).
		WillReturnError(fmt.Errorf("error"))

	films, err = repo.GetFilmsBySearch(searchStr)

	assert.Error(t, err)
	assert.Nil(t, films)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
