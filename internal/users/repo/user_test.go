package repo

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/ilyushkaaa/Filmoteka/internal/users/entity"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestRegister(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	username := "testuser"
	password := "testpassword"

	expectedUser := &entity.User{ID: 1, Username: username}

	mock.ExpectQuery("INSERT INTO users (.+) RETURNING id").
		WithArgs(username, password, "default").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedUser.ID))

	user, err := repo.Register(username, password)

	assert.NoError(t, err, "unexpected error")
	assert.NotNil(t, user, "user is nil")
	assert.Equal(t, expectedUser, user, "unexpected user returned")
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	mock.ExpectQuery("INSERT INTO users (.+) RETURNING id").
		WithArgs(username, password, "default").
		WillReturnError(fmt.Errorf("error"))

	user, err = repo.Register(username, password)

	assert.Error(t, err)
	assert.Nil(t, user)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	repo := &UserRepoPG{db: db}

	username := "testuser"

	expectedUser := &entity.User{ID: 1, Username: username}

	mock.ExpectQuery("SELECT id, username FROM users WHERE username = ?").
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(expectedUser.ID, expectedUser.Username))

	user, err := repo.GetUserByUsername(username)

	assert.NoError(t, err, "unexpected error")
	assert.NotNil(t, user, "user is nil")
	assert.Equal(t, expectedUser, user, "unexpected user returned")
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT id, username FROM users WHERE username = ?").
		WithArgs(username).
		WillReturnError(sql.ErrNoRows)

	user, err = repo.GetUserByUsername(username)

	assert.NoError(t, err, "unexpected error")
	assert.Nil(t, user, "user is nil")
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT id, username FROM users WHERE username = ?").
		WithArgs(username).
		WillReturnError(fmt.Errorf("error"))

	user, err = repo.GetUserByUsername(username)
	assert.Error(t, err)
	assert.Nil(t, user, "user is nil")
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	repo := &UserRepoPG{db: db}

	userID := uint64(1)

	expectedRole := "user"

	mock.ExpectQuery("SELECT role FROM users WHERE id = ?").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow(expectedRole))

	role, err := repo.GetUserRole(userID)
	assert.NoError(t, err, "unexpected error")
	assert.NotEmpty(t, role, "user role is empty")
	assert.Equal(t, expectedRole, role, "unexpected user role returned")
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT role FROM users WHERE id = ?").
		WithArgs(userID).
		WillReturnError(fmt.Errorf("error"))

	role, err = repo.GetUserRole(userID)
	assert.Error(t, err)
	assert.Empty(t, role)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT role FROM users WHERE id = ?").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	role, err = repo.GetUserRole(userID)
	assert.NoError(t, err)
	assert.Empty(t, role)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
