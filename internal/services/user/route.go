package user

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/umeh-promise/ecommerce/utils"
)

type Handler struct {
	store UserStore
}

func NewHandler(store UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoute() chi.Router {
	router := chi.NewRouter()

	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.registerUser)
		r.Post("/login", h.loginUser)
		r.Route("/user", func(r chi.Router) {
			r.Use(h.AuthTokenMiddleware)
			r.Get("/", h.getUser)
			r.Put("/", h.updateUser)
			r.Put("/change-password", h.changePassword)
		})
	})

	return router
}

func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := utils.ParseJSON(w, r, &payload); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	user := &User{
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Email:       payload.Email,
		Password:    hashedPassword,
		PhoneNumber: payload.PhoneNumber,
	}

	if err := h.store.CreateUser(ctx, user); err != nil {
		switch err {
		case utils.ErrorDuplicateEmail:
			utils.BadRequestError(w, r, fmt.Errorf("user with email (%s) already exists", payload.Email))
		case utils.ErrorDuplicatePhoneNumber:
			utils.BadRequestError(w, r, fmt.Errorf("user with phone number (%s) already exists", payload.PhoneNumber))
		default:
			utils.BadRequestError(w, r, err)
		}
		return
	}

	if err := utils.JSONResponse(w, http.StatusCreated, &UserResponse{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserPayload
	if err := utils.ParseJSON(w, r, &payload); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		fmt.Println(errors)
		utils.BadRequestError(w, r, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	user, err := h.store.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		utils.UnAuthorizedRequestError(w, r, fmt.Errorf("invalid email or password"))
		return
	}

	if err := utils.ComparePasswords(user.Password, payload.Password); err != nil {
		utils.UnAuthorizedRequestError(w, r, fmt.Errorf("invalid email or password"))
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	type userWithToken struct {
		User  UserResponse `json:"user"`
		Token string       `json:"token"`
	}

	userResponse := &userWithToken{
		User: UserResponse{
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		},
		Token: token,
	}

	if err := utils.JSONResponse(w, http.StatusOK, userResponse); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)

	if err := utils.JSONResponse(w, http.StatusOK, user); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	var payload UpdateUserPayload

	if err := utils.ParseJSON(w, r, &payload); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.BadRequestError(w, r, fmt.Errorf("invalid payloads, %v", errors))
		return
	}

	user := GetUserFromContext(r)

	utils.AssignIfNotNil(&user.FirstName, payload.FirstName)
	utils.AssignIfNotNil(&user.LastName, payload.LastName)
	utils.AssignIfNotNil(&user.PhoneNumber, payload.PhoneNumber)
	utils.AssignIfNotNil(&user.DOB, payload.DOB)
	utils.AssignIfNotNil(&user.Gender, payload.Gender)
	utils.AssignIfNotNil(&user.ProfilePicture, payload.ProfilePicture)

	if err := utils.JSONResponse(w, http.StatusOK, user); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

}

func (h *Handler) changePassword(w http.ResponseWriter, r *http.Request) {
	var payload ChangePasswordPayload

	user := GetUserFromContext(r)

	if err := utils.ParseJSON(w, r, &payload); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}
	if err := utils.Validator.Struct(payload); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := utils.ComparePasswords(user.Password, payload.OldPassword); err != nil {
		utils.BadRequestError(w, r, fmt.Errorf("incorrect password"))
		return
	}

	hashedPassword, err := utils.HashPassword(payload.NewPassword)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	user.Password = hashedPassword

	if err := h.store.ChangePassword(r.Context(), user); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, nil); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}
