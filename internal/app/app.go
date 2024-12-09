package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // Import the pgx driver
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
	"log"
	"musicadviser/internal/recommends"
	fridgeStore "musicadviser/internal/recommends/postgres"
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
	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return err
	}

	store := fridgeStore.NewStorage(db)

	service := recommends.NewAppService(store)
	handler := recommends.NewHandler(a.router, service)
	handler.Register()

	// shelfService := shelf.NewAppService(store)

	return nil
}

func (a *App) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	errs, ctx := errgroup.WithContext(ctx)

	log.Println("starting web server on port %s", a.config.Port)

	errs.Go(func() error {
		if err := a.http.ListenAndServe(); err != nil {
			return fmt.Errorf("listen and serve error: %w", err)
		}
		return nil
	})

	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully")

	// Perform application shutdown with a maximum timeout of 5 seconds.
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.http.Shutdown(timeoutCtx); err != nil {
		log.Println(err.Error())
	}

	return nil
}
