package usecase

import (
	"github.com/ilyushkaaa/Filmoteka/internal/actors/entity"
	"github.com/ilyushkaaa/Filmoteka/internal/actors/repo"
	"github.com/ilyushkaaa/Filmoteka/internal/dto"
)

//go:generate mockgen -source=actor.go -destination=actor_mock.go -package=usecase ActorUseCase
type ActorUseCase interface {
	GetActorByID(actorID uint64) (*dto.ActorWithFilms, error)
	GetActors() ([]dto.ActorWithFilms, error)
	AddActor(actor entity.Actor) (*entity.Actor, error)
	UpdateActor(actor entity.Actor) error
	DeleteActor(ID uint64) error
}

type ActorUseCaseApp struct {
	actorRepo repo.ActorRepo
}

func NewActorUseCase(actorRepo repo.ActorRepo) *ActorUseCaseApp {
	return &ActorUseCaseApp{
		actorRepo: actorRepo,
	}
}
func (r *ActorUseCaseApp) GetActorByID(actorID uint64) (*dto.ActorWithFilms, error) {
	actor, err := r.actorRepo.GetActorByID(actorID)
	if err != nil {
		return nil, err
	}
	if actor == nil {
		return nil, ErrActorNotFound
	}
	return actor, nil
}

func (r *ActorUseCaseApp) GetActors() ([]dto.ActorWithFilms, error) {
	actors, err := r.actorRepo.GetActors()
	if err != nil {
		return nil, err
	}
	return actors, nil
}

func (r *ActorUseCaseApp) AddActor(actor entity.Actor) (*entity.Actor, error) {
	actorID, err := r.actorRepo.AddActor(actor)
	if err != nil {
		return nil, err
	}
	actor.ID = actorID
	return &actor, nil
}

func (r *ActorUseCaseApp) UpdateActor(actor entity.Actor) error {
	wasUpdated, err := r.actorRepo.UpdateActor(actor)
	if err != nil {
		return err
	}
	if !wasUpdated {
		return ErrActorNotFound
	}
	return nil
}

func (r *ActorUseCaseApp) DeleteActor(ID uint64) error {
	wasDeleted, err := r.actorRepo.DeleteActor(ID)
	if err != nil {
		return err
	}
	if !wasDeleted {
		return ErrActorNotFound
	}
	return nil
}
