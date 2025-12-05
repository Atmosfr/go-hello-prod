package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type HelloResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hello := HelloResponse{"hello from prod-ready service", time.Now()}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(hello); err != nil {
		slog.Debug("client disconnected during response", "error", err)
	}
}
