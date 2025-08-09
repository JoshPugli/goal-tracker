package goals

import (
	"encoding/json"
	"net/http"
)

type Handlers struct {
	Store *Store
}

func NewHandlers(store *Store) *Handlers {
	return &Handlers{Store: store}
}

func (h *Handlers) HandleListGoals(w http.ResponseWriter, r *http.Request) {
	goals := h.Store.ListGoals()
	writeJSON(w, http.StatusOK, goals)
}

func (h *Handlers) HandleToday(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.Store.TodayState())
}

func (h *Handlers) HandleStats(w http.ResponseWriter, r *http.Request) {
	window := r.URL.Query().Get("window")
	if window == "" {
		window = "day"
	}
	completed, total := h.Store.Stats(window)
	writeJSON(w, http.StatusOK, Stats{Window: window, Completed: completed, Total: total})
}

// POST /api/goals/{id}/complete -> mark complete today
// DELETE /api/goals/{id}/complete -> unmark complete today
func (h *Handlers) HandleToggleComplete(w http.ResponseWriter, r *http.Request) {
	id := lastPath(r.URL.Path)
	switch r.Method {
	case http.MethodPost:
		h.Store.ToggleCompleteToday(id, true)
		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		h.Store.ToggleCompleteToday(id, false)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func lastPath(p string) string {
	// expects /api/goals/{id}/complete; take segment before last
	// naive but fine for now
	if len(p) == 0 {
		return ""
	}
	// trim trailing slash
	if p[len(p)-1] == '/' {
		p = p[:len(p)-1]
	}
	// find '/complete'
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' {
			// segment after previous '/'
			// find prev slash
			prev := i - 1
			for prev >= 0 && p[prev] != '/' {
				prev--
			}
			return p[prev+1 : i]
		}
	}
	return ""
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
