package products

import (
	"context"
	"database/sql"
	"errors"

	uuid "github.com/satori/go.uuid"
	"github.com/umeh-promise/ecommerce/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateProduct(ctx context.Context, product *Product) error {
	query := `
		INSERT INTO products
			(id,user_id,discount,name,price,description,image) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, version, created_at, updated_at
		
	`

	product.ID = uuid.NewV4().String()

	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, product.ID, product.UserID, product.Discount, product.Name, product.Price, product.Description, product.Image).Scan(
		&product.ID,
		&product.Version,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetAllProduct(ctx context.Context) ([]Product, error) {

	query := `SELECT  
		id, user_id, name, price, description, discount, image, version, created_at, updated_at
		FROM products 
	`

	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		product := Product{}
		err := rows.Scan(
			&product.ID,
			&product.UserID,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.Discount,
			&product.Image,
			&product.Version,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (s *Store) GetPostByID(ctx context.Context, id string) (*Product, error) {

	var product Product

	query := `SELECT id,user_id,discount,name,price,description,image, version, created_at, updated_at FROM products
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, id).Scan(&product.ID,
		&product.UserID,
		&product.Discount,
		&product.Name,
		&product.Price,
		&product.Description,
		&product.Image,
		&product.Version,
		&product.CreatedAt,
		&product.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, utils.ErrorNotFound
		default:
			return nil, err
		}
	}

	return &product, nil
}

func (s *Store) UpdateProduct(ctx context.Context, product *Product) error {

	query := `UPDATE products 
	SET name = $1, description = $2, image = $3, price = $4, discount=$5,  version = version + 1
	WHERE id = $6 AND version = $7
	RETURNING version
`

	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, product.Name, product.Description, product.Image, product.Price, product.Discount, product.ID, product.Version).Scan(
		&product.Version,
	)

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

func (s *Store) DeleteProduct(ctx context.Context, id string) error {
	query := `DELETE FROM products
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
