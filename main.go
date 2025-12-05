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

	"github.com/Atmosfr/go-hello-prod/internal/middleware"
	"github.com/joho/godotenv"
)

type HelloResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hello := HelloResponse{"hello from prod-ready service", time.Now()}
	res, err := json.Marshal(hello)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	var level slog.Level
	switch logLevel {
	case "info":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	slog.SetLogLoggerLevel(level)

	slog.Info("configuration loaded", "port", port, "log_level", logLevel)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloHandler)

	handler := middleware.Logging(mux)
	srv := &http.Server{Addr: ":" + port, Handler: handler}

	// separate goroutine for server
	go func() {
		slog.Info("server starting", "addr", ":"+port)
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
