package products

import "context"

type Product struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Discount    string `json:"discount"`
	Name        string `json:"name"`
	Price       string `json:"price"`
	Description string `json:"discription"`
	Image       string `json:"image"`
	Version     string `json:"-"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
}

type ProductStore interface {
	CreateProduct(context.Context, *Product) error
	GetAllProduct(context.Context) ([]Product, error)
	UpdateProduct(context.Context, *Product) error
	DeleteProduct(context.Context, string) error
	GetPostByID(context.Context, string) (*Product, error)
}

type ProductPayload struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"required,min=2"`
	Image       string `json:"image" validate:"required"`
	Price       string `json:"price" validate:"required"`
	Discount    string `json:"discount"`
}
