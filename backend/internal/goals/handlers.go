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

// HandleListGoals godoc
// @Summary List goals
// @Description Get the catalog of goals for the current user
// @Tags goals
// @Produce json
// @Success 200 {array} Goal
// @Router /api/goals [get]
func (h *Handlers) HandleListGoals(w http.ResponseWriter, r *http.Request) {
	userID := "demo" // TODO: derive from auth context
	goals := h.Store.ListGoals(userID)
	writeJSON(w, http.StatusOK, goals)
}

// HandleToday godoc
// @Summary Today's completion state
// @Description For each goal, whether it's completed today
// @Tags goals
// @Produce json
// @Success 200 {array} TodayState
// @Router /api/goals/today [get]
func (h *Handlers) HandleToday(w http.ResponseWriter, r *http.Request) {
	userID := "demo"
	writeJSON(w, http.StatusOK, h.Store.TodayState(userID))
}

// HandleStats godoc
// @Summary Stats by window
// @Description Completion stats over a time window
// @Tags goals
// @Produce json
// @Param window query string false "Time window" Enums(day,week,month) default(day)
// @Success 200 {object} Stats
// @Router /api/stats [get]
func (h *Handlers) HandleStats(w http.ResponseWriter, r *http.Request) {
	window := r.URL.Query().Get("window")
	if window == "" {
		window = "day"
	}
	userID := "demo"
	completed, total := h.Store.Stats(userID, window)
	writeJSON(w, http.StatusOK, Stats{Window: window, Completed: completed, Total: total})
}

// POST /api/goals/{id}/complete -> mark complete today
// DELETE /api/goals/{id}/complete -> unmark complete today
// HandleToggleComplete godoc
// @Summary Toggle completion for today
// @Description Mark or unmark a goal as completed for today
// @Tags goals
// @Param id path string true "Goal ID"
// @Success 204 {string} string "No Content"
// @Router /api/goals/{id}/complete [post]
// @Router /api/goals/{id}/complete [delete]
func (h *Handlers) HandleToggleComplete(w http.ResponseWriter, r *http.Request) {
	id := lastPath(r.URL.Path)
	userID := "demo"
	switch r.Method {
	case http.MethodPost:
		h.Store.ToggleCompleteToday(userID, id, true)
		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		h.Store.ToggleCompleteToday(userID, id, false)
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
