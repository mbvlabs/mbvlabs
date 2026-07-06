package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	mailclients "mbvlabs/clients/email"
	"mbvlabs/config"
	"mbvlabs/controllers"
	"mbvlabs/database"
	"mbvlabs/email"
	"mbvlabs/internal/inertia"
	"mbvlabs/internal/server"
	"mbvlabs/queue"
	"mbvlabs/router"
	"mbvlabs/services"
	"mbvlabs/telemetry"

	"go.uber.org/fx"
)

var appVersion string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := inertia.Init("views/root.go.html"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize inertia: %s\n", err)
		os.Exit(1)
	}
	app := fx.New(
		fx.Provide(
			func() context.Context { return ctx },
			func(cfg config.Config) (email.TransactionalSender, email.MarketingSender) {
				if config.Env == server.ProdEnvironment {
					log.Fatal("provide real email sender")
				}

				return mailclients.NewMailpit(cfg), mailclients.NewMailpit(cfg)
			},
		),

		config.Module,
		database.Module,
		telemetry.Module,
		queue.Module,
		queue.WorkersModule,
		services.Module,
		controllers.Module,
		router.Module,

		fx.Invoke(startQueueProcessor),
		fx.Invoke(startServer),
	)

	if err := app.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Stop(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func startQueueProcessor(lc fx.Lifecycle, p queue.Processor) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := p.Start(ctx); err != nil {
					slog.Error("queue processor error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return p.Stop(ctx)
		},
	})
}

func startServer(lc fx.Lifecycle, r *router.Router, cfg config.Config, processor queue.Processor) {
	srv := server.New(
		context.Background(),
		cfg.App.Host,
		cfg.App.Port,
		config.Env,
		r.Handler,
		[]server.Shutdowner{processor},
	)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			slog.InfoContext(
				context.Background(),
				"starting server",
				"host",
				cfg.App.Host,
				"port",
				cfg.App.Port,
			)
			go func() {
				if err := srv.Start(context.Background(), config.Env); err != nil {
					slog.Error("server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.InfoContext(ctx, "initiating graceful shutdown")
			for _, shutdowner := range srv.Shutdowners {
				if err := shutdowner.Shutdown(ctx); err != nil {
					return fmt.Errorf("component shutdown error (%T): %w", shutdowner, err)
				}
			}
			return nil
		},
	})
}
