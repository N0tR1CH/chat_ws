package main

import (
	"context"
	"embed"
	"errors"
	"html/template"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

//go:embed web/templates/*
var templatesFS embed.FS

func run(
	ctx context.Context,
	stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	config := Config{Host: "", Port: "3000"}

	parsedFS, err := template.ParseFS(templatesFS, "web/templates/*.html")
	templates := template.Must(parsedFS, err)
	if err != nil {
		slog.Error("error parsing html templates", "err", err.Error())
	}

	srv := NewServer(
		NewSlogLogger(stderr),
		templates,
	)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}
	go func() {
		slog.Info("listening and serving", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			slog.Error("error listening and serving", "err", err.Error())
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("error shutting down http server", "err", err.Error())
		}
	}()
	wg.Wait()

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(
		ctx,
		os.Stderr,
	); err != nil {
		slog.Error("error while running server", "err", err.Error())
		os.Exit(1)
	}
}
