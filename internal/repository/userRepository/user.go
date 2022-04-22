package userRepository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sword-challenge/internal/model"
)

type UserQueries interface {
	AuthenticateUser(id int) (string, error)
	GetUserByToken(token string) (*model.User, error)
	GetUsersByRole(role string) ([]model.User, error)
}

type User struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) UserQueries {
	return &User{db}
}

func (u *User) AuthenticateUser(id int) (string, error) {
	token, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	_, err = u.db.Exec(
		"INSERT INTO tokens (uuid, user_id, created_date) VALUES (?, (SELECT id FROM users WHERE id = ?), CURRENT_TIME);", token, id)
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

func (u *User) GetUserByToken(token string) (*model.User, error) {
	user := &model.User{}

	err := u.db.Get(
		user,
		"SELECT u.id, u.username, r.name as 'role.name', r.id as 'role.id' FROM users u INNER JOIN tokens t on u.id = t.user_id LEFT JOIN roles r on u.role_id = r.id WHERE t.uuid = ?;",
		token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) GetUsersByRole(role string) ([]model.User, error) {
	var users []model.User

	err := u.db.Select(
		&users,
		"SELECT u.id, u.username FROM users u INNER JOIN roles r on u.role_id = r.id WHERE r.name = ?;",
		role)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) GetUsersByName(role string) ([]model.User, error) {
	var users []model.User

	err := u.db.Select(
		&users,
		"SELECT u.id, u.username FROM users u INNER JOIN roles r on u.role_id = r.id WHERE r.name = ?;",
		role)
	if err != nil {
		return nil, err
	}

	return users, nil
}
