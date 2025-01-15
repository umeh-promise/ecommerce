package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/umeh-promise/ecommerce/utils"
	"go.uber.org/zap"
)

type APIServer struct {
	addr string
	db   *sqlx.DB
}

func NewAPIServer(addr string, db *sqlx.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}

// routerGroups ...chi.Router
func (s *APIServer) mount() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{utils.GetString("CORS_ALLOWED_ORIGIN", "https://localhost:4000")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	// router.Use(app.RateLimitMiddleware)
	router.Use(middleware.Timeout(60 * time.Second))

	// for _, subRouter := range routerGroups {
	// 	router.Mount("/v1", subRouter)
	// }

	return router
}

func (s *APIServer) Run() error {

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	handler := s.mount()

	server := &http.Server{
		Addr:         s.addr,
		Handler:      handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	logger.Info("Server has started at ", s.addr)

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		s := <-quit

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		ctx, cancel := context.WithTimeout(context.Background(), utils.QueryTimeout)
		defer cancel()

		logger.Info("Server signal ", s.String(), "caught")
		shutdown <- server.Shutdown(ctx)
	}()

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err = <-shutdown; err != nil {
		return err
	}

	logger.Info("Server existed", "addr ", s.addr)

	return nil
}
