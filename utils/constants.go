package utils

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const QueryTimeout = 5 * time.Second

var Logger *zap.SugaredLogger
var Validator *validator.Validate

func ExtendContextDuration(ctx context.Context) context.Context {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	return ctx
}

func init() {
	Logger = zap.Must(zap.NewProduction()).Sugar()
	defer Logger.Sync()

	Validator = validator.New(validator.WithRequiredStructEnabled())
}
