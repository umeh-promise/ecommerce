package user

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/umeh-promise/ecommerce/utils"
)

type userKey string

var userCtx userKey = "user"

func (middleware *Handler) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.UnAuthorizedRequestError(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnAuthorizedRequestError(w, r, fmt.Errorf("authorization header is malfarmed"))
			return
		}

		token := parts[1]
		jwtToken, err := utils.ValidateToken(token)
		if err != nil {
			utils.UnAuthorizedRequestError(w, r, err)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)
		userID := claims["sub"].(string)

		ctx := r.Context()

		user, err := middleware.store.GetUserByID(ctx, userID)
		if err != nil {
			utils.UnAuthorizedRequestError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(r *http.Request) *User {
	return r.Context().Value(userCtx).(*User)
}
