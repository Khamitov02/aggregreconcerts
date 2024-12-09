package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"musicadviser/internal/concerts"

	"golang.org/x/sync/errgroup"
	"log"

	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	config *Config
	router *chi.Mux
	http   *http.Server
}

func New(ctx context.Context, config *Config) (*App, error) {
	r := chi.NewRouter()
	return &App{
		config: config,
		router: r,
		http: &http.Server{
			Addr:              config.Host + ":" + config.Port,
			Handler:           r,
			ReadTimeout:       0,
			ReadHeaderTimeout: 0,
			WriteTimeout:      0,
			IdleTimeout:       0,
			MaxHeaderBytes:    0,
		},
	}, nil
}

func (a *App) Setup(ctx context.Context, dsn string) error {
	// Directly create the service and handler without storage
	service := concerts.NewAppService(nil) // Pass nil or a mock if needed
	handler := concerts.NewHandler(a.router, service)
	handler.Register()

	return nil
}

func (a *App) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	errs, ctx := errgroup.WithContext(ctx)

	log.Printf("starting web server on port %s", a.config.Port)

	errs.Go(func() error {
		if err := a.http.ListenAndServe(); err != nil {
			return fmt.Errorf("listen and serve error: %w", err)
		}
		return nil
	})

	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.http.Shutdown(timeoutCtx); err != nil {
		log.Println(err.Error())
	}

	return nil
}
