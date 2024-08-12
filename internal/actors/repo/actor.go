package repo

import (
	"database/sql"
	"errors"

	entityActor "github.com/ilyushkaaa/Filmoteka/internal/actors/entity"
	"github.com/ilyushkaaa/Filmoteka/internal/dto"
	entityFilm "github.com/ilyushkaaa/Filmoteka/internal/films/entity"
	"go.uber.org/zap"
)

//go:generate mockgen -source=actor.go -destination=actor_mock.go -package=repo ActorRepo
type ActorRepo interface {
	GetActorByID(actorID uint64) (*dto.ActorWithFilms, error)
	GetActors() ([]dto.ActorWithFilms, error)
	AddActor(actor entityActor.Actor) (uint64, error)
	UpdateActor(actor entityActor.Actor) (bool, error)
	DeleteActor(ID uint64) (bool, error)
}

type ActorRepoPG struct {
	db        *sql.DB
	zapLogger *zap.SugaredLogger
}

func NewActorRepo(db *sql.DB, zapLogger *zap.SugaredLogger) *ActorRepoPG {
	return &ActorRepoPG{
		db:        db,
		zapLogger: zapLogger,
	}
}
func (r *ActorRepoPG) GetActorByID(actorID uint64) (*dto.ActorWithFilms, error) {
	rows, err := r.db.Query(`
        SELECT a.id, a.name, a.surname, a.gender, a.birthday, f.id, f.name, f.description, f.date_of_release, f.rating
        FROM actors a
        LEFT JOIN film_actors fa ON a.id = fa.actor_id
        LEFT JOIN films f ON fa.film_id = f.id
        WHERE a.id = $1
    `, actorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			r.zapLogger.Errorf("error in closing query rows: %s", err)
		}
	}(rows)

	var actorWithFilms dto.ActorWithFilms
	actorWithFilms.Films = make([]entityFilm.Film, 0)

	for rows.Next() {
		var actor entityActor.Actor
		var filmDB dto.FilmDB
		err = rows.Scan(&actor.ID, &actor.Name, &actor.Surname, &actor.Gender, &actor.Birthday, &filmDB.ID, &filmDB.Name,
			&filmDB.Description, &filmDB.DateOfRelease, &filmDB.Rating)
		if err != nil {
			return nil, err
		}

		actorWithFilms.Actor = actor
		film := filmDB.GetFilm()
		if film != nil {
			actorWithFilms.Films = append(actorWithFilms.Films, *film)
		}
	}
	if actorWithFilms.Actor.Name == "" {
		return nil, nil
	}

	return &actorWithFilms, nil
}

func (r *ActorRepoPG) GetActors() ([]dto.ActorWithFilms, error) {
	rows, err := r.db.Query(`
        SELECT a.id, a.name, a.surname, a.gender, a.birthday, f.id, f.name, f.description, f.date_of_release, f.rating
        FROM actors a
        LEFT JOIN film_actors fa ON a.id = fa.actor_id
        LEFT JOIN films f ON fa.film_id = f.id
    `)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			r.zapLogger.Errorf("error in closing query rows: %s", err)
		}
	}(rows)

	actorsMap := make(map[int]dto.ActorWithFilms)
	for rows.Next() {
		var actorID int
		var actor entityActor.Actor
		var filmDB dto.FilmDB
		err = rows.Scan(&actorID, &actor.Name, &actor.Surname, &actor.Gender, &actor.Birthday, &filmDB.ID, &filmDB.Name,
			&filmDB.Description, &filmDB.DateOfRelease, &filmDB.Rating)
		if err != nil {
			return nil, err
		}

		film := filmDB.GetFilm()
		if awf, ok := actorsMap[actorID]; ok {
			if film != nil {
				awf.Films = append(awf.Films, *film)
			}
			actorsMap[actorID] = awf
		} else {
			actorFilms := make([]entityFilm.Film, 0)
			if film != nil {
				actorFilms = append(actorFilms, *film)
			}
			actorsMap[actorID] = dto.ActorWithFilms{
				Actor: actor,
				Films: actorFilms,
			}
		}
	}

	actorsWithFilms := make([]dto.ActorWithFilms, 0, len(actorsMap))
	for _, awf := range actorsMap {
		actorsWithFilms = append(actorsWithFilms, awf)
	}

	return actorsWithFilms, nil
}

func (r *ActorRepoPG) AddActor(actor entityActor.Actor) (uint64, error) {
	var actorID uint64
	err := r.db.
		QueryRow("INSERT INTO actors (name, surname, gender, birthday) VALUES ($1, $2, $3, $4) RETURNING id",
			actor.Name, actor.Surname, actor.Gender, actor.Birthday).
		Scan(&actorID)
	return actorID, err
}

func (r *ActorRepoPG) UpdateActor(actor entityActor.Actor) (bool, error) {
	_, err := r.db.Exec("UPDATE actors SET name = $1, surname = $2, gender = $3, birthday = $4 WHERE id = $5",
		actor.Name, actor.Surname, actor.Gender, actor.Birthday, actor.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *ActorRepoPG) DeleteActor(ID uint64) (bool, error) {
	result, err := r.db.Exec("DELETE FROM actors WHERE id = $1", ID)
	if err != nil {
		return false, err
	}
	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsDeleted > 0, nil
}
