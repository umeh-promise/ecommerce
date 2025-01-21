package user

import (
	"context"
	"database/sql"

	uuid "github.com/satori/go.uuid"
	"github.com/umeh-promise/ecommerce/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users(id, first_name, last_name, email, password, phone_number, dob, gender, profile_picture) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, version, created_at, updated_at
	`
	user.ID = uuid.NewV4().String()

	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query,
		user.ID, user.FirstName, user.LastName,
		user.Email, user.Password,
		user.PhoneNumber, user.DOB,
		user.Gender, user.ProfilePicture).Scan(
		&user.ID,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return utils.ErrorDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_phone_number_key"`:
			return utils.ErrorDuplicatePhoneNumber
		default:
			return err
		}
	}

	return nil
}

func (s *Store) GetUserByID(ctx context.Context, userID string) (*User, error) {
	var user User

	query := `
		SELECT id, first_name, last_name, email, phone_number, dob, gender, profile_picture, password, version FROM users
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.FirstName,
		&user.LastName, &user.Email, &user.PhoneNumber,
		&user.DOB, &user.Gender,
		&user.ProfilePicture, &user.Password, &user.Version,
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
		SELECT id, first_name, last_name, email, password, phone_number, version FROM users
		WHERE email = $1
	`
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.FirstName,
		&user.LastName, &user.Email,
		&user.Password, &user.PhoneNumber,
		&user.Version,
	)
	if err != nil {
		return &User{}, err
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

	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()
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

func (s *Store) DeleteUser(ctx context.Context, userID string) error {
	query := `
	DELETE FROM users 
	WHERE id = $1
`
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

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

func (s *Store) ChangePassword(ctx context.Context, user *User) error {
	query := `
	UPDATE users 
	SET password = $1
	WHERE id = $2
	RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, user.Password, user.ID).Scan(&user.Version)
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
