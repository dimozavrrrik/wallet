package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimozavrrrik/wallet/internal/config"
	"github.com/dimozavrrrik/wallet/internal/handler"
	repo "github.com/dimozavrrrik/wallet/internal/repository"
	"github.com/dimozavrrrik/wallet/internal/service"
	"github.com/dimozavrrrik/wallet/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("init pool: %w", err)
	}
	defer pool.Close()

	pingCtx, pingCancel := context.WithTimeout(ctx, time.Second*5)
	defer pingCancel()
	if err = pool.Ping(pingCtx); err != nil {
		return fmt.Errorf("ping pool: %w", err)
	}

	log.Println("connected to database")

	err = migrations.RunMigrations(ctx, pool)
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	walletRepo := repo.NewWalletRepository(pool)
	walletSvc := service.NewWalletService(walletRepo)
	walletHandler := handler.NewWalletHandler(walletSvc)

	walletHandler.Register(mux)

	srv := &http.Server{
		Addr:        cfg.HTTPAddr,
		Handler:     mux,
		ReadTimeout: time.Second * 5,
	}

	srvErr := make(chan error, 1)
	go func() {
		log.Println("server listening on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
		}
	}()

	select {
	case err := <-srvErr:
		return fmt.Errorf("http server error: %w", err)
	case <-ctx.Done():
		log.Println("shutting down server")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	if err = srv.Shutdown(shutdownCtx); err != nil {
		log.Println("shutdown error:", err)
	}

	log.Println("server stopped")

	return nil

}
