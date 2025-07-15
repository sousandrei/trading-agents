package orchestrator

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/sousandrei/trading-agents/internal/types"
)

func (o *Orchestrator) Handler() (string, http.Handler) {
	return "/analyse", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var positions []types.Position
		if err := json.NewDecoder(r.Body).Decode(&positions); err != nil {
			log.Println("failed to decode positions: ", slog.Any("err", err))
			http.Error(w, "Invalid positions format", http.StatusBadRequest)
			return
		}

		actions := []types.Action{}

		for _, position := range positions {
			res, err := o.Analyze(r.Context(), position)
			if err != nil {
				log.Println("failed to analyze: ", slog.Any("err", err))
				return
			}

			actions = append(actions, *res)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(actions); err != nil {
			log.Println("failed to encode actions: ", slog.Any("err", err))
			http.Error(w, "Failed to encode actions", http.StatusInternalServerError)
			return
		}
	})
}
