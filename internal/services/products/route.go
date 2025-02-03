package products

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/umeh-promise/ecommerce/internal/services/user"
	"github.com/umeh-promise/ecommerce/utils"
)

type Handler struct {
	store ProductStore
}

func NewHandler(store ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoute(auth *user.Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/products", func(r chi.Router) {
			r.With(auth.AuthTokenMiddleware).Post("/", h.createProduct)
			r.Get("/", h.getAllProduct)
			r.Route("/{id}", func(r chi.Router) {
				r.Use(h.ProductMiddleware)
				r.Get("/", h.getProduct)
				r.Put("/", h.updateProduct)
				r.Delete("/", h.deleteProduct)
			})
		})
	}

}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var payload ProductPayload

	if err := utils.ParseJSON(w, r, &payload); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	user := user.GetUserFromContext(r)

	ctx := r.Context()

	product := &Product{
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		UserID:      user.ID,
		Discount:    payload.Discount,
		Price:       payload.Price,
	}

	if err := h.store.CreateProduct(ctx, product); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if err := utils.JSONResponse(w, http.StatusCreated, product); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

}

func (h *Handler) getAllProduct(w http.ResponseWriter, r *http.Request) {

	products, err := h.store.GetAllProduct(r.Context())
	if err != nil {
		utils.NotFoundResponse(w, r, err)
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, products); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}

func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request) {
	product := GetProductFromMiddleware(r)

	if err := utils.JSONResponse(w, http.StatusOK, product); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}

func (h *Handler) updateProduct(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name        *string `json:"name" validate:"omitempty"`
		Description *string `json:"description" validate:"omitempty"`
		Price       *string `json:"price" validate:"omitempty"`
		Image       *string `json:"image" validate:"omitempty"`
		Discount    *string `json:"discount" validate:"omitempty"`
	}

	product := GetProductFromMiddleware(r)

	if err := utils.ParseJSON(w, r, &payload); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	utils.AssignIfNotNil(&product.Name, payload.Name)
	utils.AssignIfNotNil(&product.Description, payload.Description)
	utils.AssignIfNotNil(&product.Price, payload.Price)
	utils.AssignIfNotNil(&product.Image, payload.Image)
	utils.AssignIfNotNil(&product.Discount, payload.Discount)

	err := h.store.UpdateProduct(r.Context(), product)
	if err != nil {
		switch err {
		case utils.ErrorNotFound:
			utils.NotFoundResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, product); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

}

func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	product := GetProductFromMiddleware(r)

	if err := h.store.DeleteProduct(r.Context(), product.ID); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if err := utils.JSONResponse(w, http.StatusNoContent, nil); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}
