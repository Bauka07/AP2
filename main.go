package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Bauka07/AP2/internal/server"
	"github.com/Bauka07/AP2/internal/store"
	"github.com/Bauka07/AP2/internal/worker"
)

func main() {


	st := store.NewStore[string, string]()
	srv := server.NewServer(st)


	httpServer := &http.Server{
		Addr: ":8080",
		Handler: srv.Router(),
	}

	go worker.StartWorker(srv)

	// start server

	serverErrCh := make(chan error, 1)

	go func() {
		fmt.Println("Server starting on :8080")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
	}()

	//graceful shuwdown

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case sig := <-sigCh:
		fmt.Printf("Received signal: %v. Shutting down...", sig)
	case err := <-serverErrCh:
		fmt.Printf("Server error: %v", err)
	}
	//stopping worker
	worker.StopWorker(srv)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Graceful shutdown failed: %v", err)
		_ = httpServer.Close()
	}

	fmt.Println("Server exited")
}
