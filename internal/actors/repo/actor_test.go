package repo

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	entityActor "github.com/ilyushkaaa/Filmoteka/internal/actors/entity"
	"github.com/ilyushkaaa/Filmoteka/internal/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetActors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can not create mock")
	}
	defer db.Close()
	testRepo := NewActorRepo(db, zap.NewNop().Sugar())

	var nilActors []dto.ActorWithFilms

	mock.ExpectQuery(`
        SELECT a.id, a.name, a.surname, a.gender, a.birthday, f.id, f.name, f.description, f.date_of_release, f.rating
        FROM actors a
        LEFT JOIN film_actors fa ON a.id = fa.actor_id
        LEFT JOIN films f ON fa.film_id = f.id
    `).WillReturnError(fmt.Errorf("error"))
	actor, err := testRepo.GetActors()
	if errExp := mock.ExpectationsWereMet(); errExp != nil {
		t.Errorf("there were unfulfilled expectations: %s", errExp)
		return
	}
	assert.NotEqual(t, nil, err)
	assert.Equal(t, nilActors, actor)

	actorRows := sqlmock.NewRows([]string{"id", "name", "surname", "gender", "birthday", "f_id", "f_name", "f_description", "f_date_of_release", "f_rating"}).
		AddRow(1, "John", "Doe", "Male", time.Time{}.Add(time.Hour), 1, "Film 1", "Description 1", time.Time{}.Add(time.Hour), 8.0).
		AddRow(2, "Jane", "Smith", "Female", time.Time{}.Add(time.Hour), 2, "Film 2", "Description 2", time.Time{}.Add(time.Hour), 7.5)

	mock.ExpectQuery(`
        SELECT a.id, a.name, a.surname, a.gender, a.birthday, f.id, f.name, f.description, f.date_of_release, f.rating
        FROM actors a
        LEFT JOIN film_actors fa ON a.id = fa.actor_id
        LEFT JOIN films f ON fa.film_id = f.id
    `).WillReturnRows(actorRows)

	actors, err := testRepo.GetActors()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(actors))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetActorByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can not create mock")
	}
	defer db.Close()
	testRepo := NewActorRepo(db, zap.NewNop().Sugar())

	var id uint64 = 1
	var nilActorWithFilms *dto.ActorWithFilms

	mock.
		ExpectQuery(`
        SELECT a.id, a.name, a.surname, a.gender, a.birthday, f.id, f.name, f.description, f.date_of_release, f.rating
        FROM actors a
        LEFT JOIN film_actors fa ON a.id = fa.actor_id
        LEFT JOIN films f ON fa.film_id = f.id
        WHERE a.id
    `).WithArgs(id).WillReturnError(fmt.Errorf("db_error"))

	actor, err := testRepo.GetActorByID(id)
	if errExp := mock.ExpectationsWereMet(); errExp != nil {
		t.Errorf("there were unfulfilled expectations: %s", errExp)
		return
	}
	assert.NotEqual(t, nil, err)
	assert.Equal(t, nilActorWithFilms, actor)

	mock.
		ExpectQuery(`
        SELECT a.id, a.name, a.surname, a.gender, a.birthday, f.id, f.name, f.description, f.date_of_release, f.rating
        FROM actors a
        LEFT JOIN film_actors fa ON a.id = fa.actor_id
        LEFT JOIN films f ON fa.film_id = f.id
        WHERE a.id
    `).WithArgs(id).WillReturnError(sql.ErrNoRows)

	actor, err = testRepo.GetActorByID(id)
	if errExp := mock.ExpectationsWereMet(); errExp != nil {
		t.Errorf("there were unfulfilled expectations: %s", errExp)
		return
	}
	assert.Equal(t, nil, err)
	assert.Equal(t, nilActorWithFilms, actor)

	var expectedActorID uint64 = 1
	var expectedActorName = "John"
	var expectedFilmID uint64 = 1
	var expectedFilmName = "Film 1"
	mock.ExpectQuery(`
        SELECT a.id, a.name, a.surname, a.gender, a.birthday, f.id, f.name, f.description, f.date_of_release, f.rating
        FROM actors a
        LEFT JOIN film_actors fa ON a.id = fa.actor_id
        LEFT JOIN films f ON fa.film_id = f.id
        WHERE a.id
    `).WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"a.id", "a.name", "a.surname", "a.gender", "a.birthday", "f.id", "f.name", "f.description", "f.date_of_release", "f.rating"}).
			AddRow(expectedActorID, expectedActorName, "Doe", "male", time.Time{}.Add(time.Hour), expectedFilmID, expectedFilmName, "Film Description", time.Time{}.Add(time.Hour), 8.0))

	actor, err = testRepo.GetActorByID(id)

	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, actor)
	assert.Equal(t, expectedActorID, actor.Actor.ID)
	assert.Equal(t, expectedActorName, actor.Actor.Name)
	assert.Equal(t, 1, len(actor.Films))

	err = mock.ExpectationsWereMet()
	assert.Equal(t, nil, err)
}

func TestAddActor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	testRepo := NewActorRepo(db, zap.NewNop().Sugar())

	expectedActorID := uint64(123)

	mock.ExpectQuery("INSERT INTO actors (.+) RETURNING id").
		WithArgs("John", "Doe", "Male", time.Time{}.Add(time.Hour)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedActorID))

	actor := entityActor.Actor{Name: "John", Surname: "Doe", Gender: "Male", Birthday: time.Time{}.Add(time.Hour)}
	actorID, err := testRepo.AddActor(actor)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedActorID, actorID)
	err = mock.ExpectationsWereMet()
	assert.Equal(t, nil, err)

}

func TestUpdateActor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	testRepo := NewActorRepo(db, zap.NewNop().Sugar())

	mock.ExpectExec("UPDATE actors SET (.+) WHERE id = (.+)").
		WithArgs("John", "Doe", "Male", time.Time{}.Add(1), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	actor := entityActor.Actor{ID: 1, Name: "John", Surname: "Doe", Gender: "Male", Birthday: time.Time{}.Add(1)}
	success, err := testRepo.UpdateActor(actor)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, success)
	err = mock.ExpectationsWereMet()
	assert.Equal(t, nil, err)

	mock.ExpectExec("UPDATE actors SET (.+) WHERE id = (.+)").
		WithArgs("John", "Doe", "Male", time.Time{}.Add(1), 1).
		WillReturnError(fmt.Errorf("error"))

	success, err = testRepo.UpdateActor(actor)

	assert.NotEqual(t, nil, err)
	assert.Equal(t, false, success)
	err = mock.ExpectationsWereMet()
	assert.Equal(t, nil, err)
}
