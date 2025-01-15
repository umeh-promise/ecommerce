package user

import (
	"context"
	"database/sql"

	"github.com/umeh-promise/ecommerce/utils"
)

type Store struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users(first_name, last_name, email, password, phone_number, dob, gender, profile_picture) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING version, created_at, updated_at
	`

	ctx = utils.ExtendContextDuration(ctx)

	err := s.db.QueryRowContext(ctx, query,
		user.FirstName, user.LastName,
		user.Email, user.Password,
		user.PhoneNumber, user.DOB,
		user.Gender, user.ProfilePicture).Scan(
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByID(ctx context.Context, userID string) (*User, error) {
	var user User

	query := `
		SELECT id, first_name, last_name, email, password, phone_number, dob, gender, profile_picture, version FROM users
		WHERE id = $1
	`
	ctx = utils.ExtendContextDuration(ctx)

	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.FirstName,
		&user.LastName, &user.Email,
		&user.Password, &user.PhoneNumber,
		&user.DOB, &user.Gender,
		&user.ProfilePicture, &user.Version,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, utils.ErrorNotFound
		default:
			return nil, err

		}
	}

	return &user, nil

}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	query := `
		SELECT id, first_name, last_name, email, password, phone_number, dob, gender, profile_picture, version FROM users
		WHERE email = $1
	`
	ctx = utils.ExtendContextDuration(ctx)

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.FirstName,
		&user.LastName, &user.Email,
		&user.Password, &user.PhoneNumber,
		&user.DOB, &user.Gender,
		&user.ProfilePicture, &user.Version,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, utils.ErrorNotFound
		default:
			return nil, err

		}
	}

	return &user, nil

}

func (s *Store) UpdateUser(ctx context.Context, user *User) error {
	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, dob = $3, gender = $4, profile_picture = $5
		WHERE id = $6
		RETURNING version
	`

	ctx = utils.ExtendContextDuration(ctx)
	err := s.db.QueryRowContext(ctx, query, user.ID).Scan(&user.Version)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return utils.ErrorNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *Store) DeleteUser(ctx context.Context, userID int64) error {
	query := `
	DELETE FROM users 
	WHERE id = $1
`
	ctx = utils.ExtendContextDuration(ctx)

	rows, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	row, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if row <= 0 {
		return utils.ErrorNotFound
	}

	return nil
}
