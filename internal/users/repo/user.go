package repo

import (
	"database/sql"
	"errors"

	"github.com/ilyushkaaa/Filmoteka/internal/users/entity"
)

//go:generate mockgen -source=user.go -destination=user_mock.go -package=repo UserRepo
type UserRepo interface {
	Login(username, password string) (*entity.User, error)
	Register(username, password string) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	GetUserRole(userID uint64) (string, error)
}

type UserRepoPG struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepoPG {
	return &UserRepoPG{
		db: db,
	}
}

func (u *UserRepoPG) Login(username, password string) (*entity.User, error) {
	foundUser := &entity.User{}
	err := u.db.
		QueryRow("SELECT id, username FROM users WHERE username = $1 AND password = $2", username, password).
		Scan(&foundUser.ID, &foundUser.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return foundUser, nil
}

func (u *UserRepoPG) Register(username, password string) (*entity.User, error) {
	var userID uint64
	err := u.db.
		QueryRow("INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING id", username, password, "default").
		Scan(&userID)
	if err != nil {
		return nil, err
	}
	return &entity.User{
		ID:       userID,
		Username: username,
	}, nil
}

func (u *UserRepoPG) GetUserByUsername(username string) (*entity.User, error) {
	foundUser := &entity.User{}
	err := u.db.
		QueryRow("SELECT id, username FROM users WHERE username = $1", username).
		Scan(&foundUser.ID, &foundUser.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return foundUser, nil
}

func (u *UserRepoPG) GetUserRole(userID uint64) (string, error) {
	var userRole string
	err := u.db.
		QueryRow("SELECT role FROM users WHERE id = $1", userID).
		Scan(&userRole)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return userRole, nil
}
