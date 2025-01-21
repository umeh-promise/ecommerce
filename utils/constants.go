package utils

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const QueryTimeout = 5 * time.Second

var Logger *zap.SugaredLogger
var Validator *validator.Validate

var (
	tokenExp    = time.Hour * 24
	tokenIssuer = "ecommerce"
	authSecret  = GetString("AUTH_SECRET", "basic")
)

func init() {
	Logger = zap.Must(zap.NewProduction()).Sugar()
	defer Logger.Sync()

	Validator = validator.New(validator.WithRequiredStructEnabled())
}

func AssignIfNotNil(dest *string, src *string) {
	if src != nil {
		*dest = *src
	}
}
