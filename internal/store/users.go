package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type Password struct {
	text *string
	hash []byte
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		INSERT INTO users (first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `
		SELECT id , first_name , last_name , email , password , created_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}

	err := s.db.QueryRowContext(ctx,
		query,
		userId,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
