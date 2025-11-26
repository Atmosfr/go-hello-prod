package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HelloResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("incoming request", "method", r.Method, "path", r.URL.Path)

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hello := HelloResponse{"hello from prod-ready service", time.Now()}
	res, err := json.Marshal(hello)
	if err != nil {
		slog.Error("failed to marshal response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		slog.Error("failed to write response", "error", err)
	}
	durationMs := float64(time.Since(start).Microseconds()) / 1000
	slog.Info("request completed",
		"method", r.Method,
		"path", r.URL.Path,
		"duration_ms", durationMs,
	)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloHandler)

	srv := &http.Server{Addr: ":8080", Handler: mux}

	// separate goroutine for server
	go func() {
		slog.Info("server starting", "addr", "0.0.0.0:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server crashed", "err", err)
			os.Exit(1)
		}
	}()

	//chan for system signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "err", err)
	} else {
		slog.Info("server stopped cleanly")
	}

}
