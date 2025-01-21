package utils

import (
	"errors"
	"net/http"
)

var (
	ErrorNotFound             = errors.New("resource not found")
	ErrorInvalidID            = errors.New("invalid post id")
	ErrorDuplicateEmail       = errors.New("a user with that email already exists")
	ErrorDuplicatePhoneNumber = errors.New("duplicate phone number")
)

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Errorw("internal server error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error())

	WriteJSONError(w, http.StatusInternalServerError, []string{}, "The server encountered a problem")
}

func ForbiddenServerError(w http.ResponseWriter, r *http.Request) {
	Logger.Errorw("internal server error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", "forbidden")

	WriteJSONError(w, http.StatusForbidden, []string{}, "Forbidden")
}
func BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Errorw("bad request",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error())

	errors := []string{err.Error()}

	WriteJSONError(w, http.StatusBadRequest, errors, "validation errors")
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Errorw("not found error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error())

	WriteJSONError(w, http.StatusNotFound, []string{}, "not found")
}

func UnAuthorizedRequestError(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Errorw("unauthorized request",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error())

	errors := []string{}

	WriteJSONError(w, http.StatusUnauthorized, errors, err.Error())
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	Logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	WriteJSONError(w, http.StatusTooManyRequests, []string{}, "rate limit exceeded, retry after: "+retryAfter)
}
