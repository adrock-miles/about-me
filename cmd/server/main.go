package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/adrock-miles/about-me/internal/handlers"
	"github.com/adrock-miles/about-me/internal/platform/config"
	"github.com/adrock-miles/about-me/internal/platform/notify"
	"github.com/adrock-miles/about-me/internal/platform/ratelimit"
	"github.com/adrock-miles/about-me/internal/platform/server"
	"github.com/adrock-miles/about-me/internal/web"
)

func main() {
	root := &cobra.Command{
		Use:   "about-me",
		Short: "Personal portfolio site server",
		Long:  "Server-rendered portfolio site (Go + Datastar + Tailwind v4).",
		RunE:  runServe,
	}

	root.AddCommand(serveCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func serveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server (default)",
		RunE:  runServe,
	}
	return cmd
}

// init registers `serve` as the default subcommand so `./about-me` and
// `./about-me serve` behave the same — handy for entrypoints/PaaS.
func init() {
	cobra.EnableCommandSorting = false
}

func runServe(cmd *cobra.Command, _ []string) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	mailer := notify.NewMailer(cfg.Email, logger)

	tmpl, err := handlers.ParseTemplates(web.Templates, web.AssetVersion)
	if err != nil {
		return fmt.Errorf("loading templates: %w", err)
	}
	page := handlers.NewPage(tmpl, logger)
	contactLimiter := ratelimit.New(5, time.Hour)
	contact := handlers.NewContact(tmpl, mailer, logger, contactLimiter)

	staticFS, err := fs.Sub(web.Static, "static")
	if err != nil {
		return fmt.Errorf("static fs: %w", err)
	}

	router := server.NewRouter(server.Deps{
		Page:    page,
		Contact: contact,
		Static:  staticFS,
		Logger:  logger,
	})

	return server.Run(cmd.Context(), cfg.Server, logger, router)
}
