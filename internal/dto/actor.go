package dto

import (
	"time"

	"github.com/asaskevich/govalidator"
	entityActor "github.com/ilyushkaaa/Filmoteka/internal/actors/entity"
	entityFilm "github.com/ilyushkaaa/Filmoteka/internal/films/entity"
	"github.com/ilyushkaaa/Filmoteka/pkg/validator"
)

type (
	ActorWithFilms struct {
		Actor entityActor.Actor
		Films []entityFilm.Film
	}
	ActorAdd struct {
		Name     string    `json:"name" valid:"required,length(1|40)"`
		Surname  string    `json:"surname" valid:"required,length(1|40)"`
		Gender   string    `json:"gender" valid:"required,in(male|female)"`
		Birthday time.Time `json:"birthday" valid:"required"`
	}
	ActorUpdate struct {
		ID       uint64    `json:"id" valid:"required"`
		Name     string    `json:"name" valid:"required,length(1|40)"`
		Surname  string    `json:"surname" valid:"required,length(1|40)"`
		Gender   string    `json:"gender" valid:"required,in(male|female)"`
		Birthday time.Time `json:"birthday" valid:"required"`
	}
)

func (a *ActorAdd) Validate() []string {
	_, err := govalidator.ValidateStruct(a)
	return validator.CollectErrors(err)
}

func (a *ActorUpdate) Validate() []string {
	_, err := govalidator.ValidateStruct(a)
	return validator.CollectErrors(err)
}

func (a *ActorAdd) Convert() entityActor.Actor {
	return entityActor.Actor{
		Name:     a.Name,
		Surname:  a.Surname,
		Gender:   a.Gender,
		Birthday: a.Birthday,
	}
}

func (a *ActorUpdate) Convert() entityActor.Actor {
	return entityActor.Actor{
		ID:       a.ID,
		Name:     a.Name,
		Surname:  a.Surname,
		Gender:   a.Gender,
		Birthday: a.Birthday,
	}
}
