// Package server builds the chi router and runs the HTTP listener with
// graceful shutdown.
package server

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/adrock-miles/about-me/internal/handlers"
	"github.com/adrock-miles/about-me/internal/platform/config"
	mw "github.com/adrock-miles/about-me/internal/platform/middleware"
)

// Deps is the bundle of services the router needs.
type Deps struct {
	Page    *handlers.Page
	Contact *handlers.Contact
	Static  fs.FS
	Logger  *slog.Logger
}

// NewRouter returns a fully wired chi.Router.
func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(mw.StructuredLogger(d.Logger))
	r.Use(chimw.Recoverer)
	r.Use(chimw.Compress(5))
	r.Use(mw.SecurityHeaders)

	r.Get("/health", handlers.Health)

	r.Get("/", d.Page.Index)
	r.Post("/contact", d.Contact.Submit)

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/favicon.svg", http.StatusMovedPermanently)
	})

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(d.Static))))

	return r
}

// Run starts the HTTP server and blocks until SIGINT/SIGTERM is received,
// then drains in-flight requests within cfg.ShutdownGrace before returning.
func Run(ctx context.Context, cfg config.ServerConfig, logger *slog.Logger, h http.Handler) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      h,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		logger.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("listen: %w", err)
		}
		return nil
	case sig := <-stop:
		logger.Info("shutdown signal received", "signal", sig.String())
	case <-ctx.Done():
		logger.Info("context canceled, shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownGrace)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	logger.Info("server stopped cleanly")
	return nil
}
