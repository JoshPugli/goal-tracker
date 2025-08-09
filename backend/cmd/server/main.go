package main

import (
	"fmt"
	"github.com/JoshPugli/grindhouse-api/internal/api"
	"net/http"
	"os"
	"os/signal"
	"context"
	"io"
	"syscall"
	"time"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	srv := api.NewServer()
	server := &http.Server{
		Addr:         ":8000",
		Handler:      srv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		fmt.Fprintf(w, "Server listening on port :8000\n")
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		fmt.Fprintf(w, "\nShutdown signal received...\n")
		
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		
		if err := server.Shutdown(shutdownCtx); err != nil {
			// Force shutdown
			server.Close()
			return fmt.Errorf("could not gracefully shut down server: %w", err)
		}
		
		fmt.Fprintf(w, "Server gracefully shut down\n")
	}

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}