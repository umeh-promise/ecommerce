package products

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/umeh-promise/ecommerce/utils"
)

type productKey string

var productCtx productKey = "product"

func (middleware *Handler) ProductMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		productID := chi.URLParam(r, "id")
		ctx := r.Context()

		product, err := middleware.store.GetPostByID(ctx, productID)
		if err != nil {
			switch err {
			case utils.ErrorNotFound:
				utils.NotFoundResponse(w, r, err)
			default:
				utils.InternalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, productCtx, product)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func GetProductFromMiddleware(r *http.Request) *Product {
	return r.Context().Value(productCtx).(*Product)
}
