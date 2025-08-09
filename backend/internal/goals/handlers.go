package goals

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JoshPugli/grindhouse-api/internal/auth"
)

type Handlers struct {
	goalRepo *Repository
}

func NewHandlers(goalRepo *Repository) *Handlers {
	return &Handlers{
		goalRepo: goalRepo,
	}
}

// HandleCreateGoal godoc
// @Summary Create a new goal
// @Description Create a new goal for the authenticated user
// @Tags goals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param goal body CreateGoalRequest true "Goal data"
// @Success 201 {object} Goal
// @Failure 400 {string} string "Invalid JSON or validation error"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/goals [post]
func (h *Handlers) HandleCreateGoal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req CreateGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if req.GoalType == "" {
		http.Error(w, "Goal type is required", http.StatusBadRequest)
		return
	}

	if req.GoalType != GoalTypeBoolean && req.GoalType != GoalTypeNumeric && req.GoalType != GoalTypeDuration {
		http.Error(w, "Invalid goal type", http.StatusBadRequest)
		return
	}

	goal, err := h.goalRepo.CreateGoal(userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(goal)
}

// HandleGetGoals godoc
// @Summary Get user's goals
// @Description Get all active goals for the authenticated user
// @Tags goals
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Goal
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/goals [get]
func (h *Handlers) HandleGetGoals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	goals, err := h.goalRepo.GetGoalsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(goals)
}

// HandleGetGoalsToday godoc
// @Summary Get user's goals with today's instances
// @Description Get all active goals for the authenticated user with today's daily instances
// @Tags goals
// @Produce json
// @Security BearerAuth
// @Success 200 {array} GoalWithTodayInstance
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/goals/today [get]
func (h *Handlers) HandleGetGoalsToday(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	goals, err := h.goalRepo.GetGoalsWithTodayInstances(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(goals)
}

// HandleGetGoal godoc
// @Summary Get a specific goal
// @Description Get a specific goal by ID for the authenticated user
// @Tags goals
// @Produce json
// @Security BearerAuth
// @Param id path string true "Goal ID"
// @Success 200 {object} Goal
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Goal not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/goals/{id} [get]
func (h *Handlers) HandleGetGoal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	goalID := r.URL.Path[len("/api/goals/"):]
	if goalID == "" {
		http.Error(w, "Goal ID is required", http.StatusBadRequest)
		return
	}

	goal, err := h.goalRepo.GetGoalByID(goalID, userID)
	if err != nil {
		if err.Error() == "goal not found" {
			http.Error(w, "Goal not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(goal)
}

// HandleUpdateGoal godoc
// @Summary Update a goal
// @Description Update a specific goal by ID for the authenticated user
// @Tags goals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Goal ID"
// @Param goal body UpdateGoalRequest true "Updated goal data"
// @Success 200 {object} Goal
// @Failure 400 {string} string "Invalid JSON"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Goal not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/goals/{id} [put]
func (h *Handlers) HandleUpdateGoal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	goalID := r.URL.Path[len("/api/goals/"):]
	if goalID == "" {
		http.Error(w, "Goal ID is required", http.StatusBadRequest)
		return
	}

	var req UpdateGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	goal, err := h.goalRepo.UpdateGoal(goalID, userID, req)
	if err != nil {
		if err.Error() == "goal not found" {
			http.Error(w, "Goal not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(goal)
}

// HandleDeleteGoal godoc
// @Summary Delete a goal
// @Description Soft delete a specific goal by ID for the authenticated user
// @Tags goals
// @Security BearerAuth
// @Param id path string true "Goal ID"
// @Success 204 "No Content"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Goal not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/goals/{id} [delete]
func (h *Handlers) HandleDeleteGoal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	goalID := r.URL.Path[len("/api/goals/"):]
	if goalID == "" {
		http.Error(w, "Goal ID is required", http.StatusBadRequest)
		return
	}

	err := h.goalRepo.DeleteGoal(goalID, userID)
	if err != nil {
		if err.Error() == "goal not found" {
			http.Error(w, "Goal not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleUpdateDailyInstance godoc
// @Summary Update daily goal instance
// @Description Update a daily goal instance for a specific date
// @Tags goals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param goalId path string true "Goal ID"
// @Param date query string false "Date (YYYY-MM-DD format, defaults to today)"
// @Param instance body UpdateDailyInstanceRequest true "Daily instance data"
// @Success 200 {object} DailyGoalInstance
// @Failure 400 {string} string "Invalid JSON or date format"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Goal not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/goals/{goalId}/daily [put]
func (h *Handlers) HandleUpdateDailyInstance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	goalID := r.URL.Path[len("/api/goals/"):]
	goalID = goalID[:len(goalID)-len("/daily")]
	if goalID == "" {
		http.Error(w, "Goal ID is required", http.StatusBadRequest)
		return
	}

	dateStr := r.URL.Query().Get("date")
	var date time.Time
	if dateStr == "" {
		date = time.Now()
	} else {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}

	instance, err := h.goalRepo.GetOrCreateDailyInstance(goalID, userID, date)
	if err != nil {
		if err.Error() == "goal not found" {
			http.Error(w, "Goal not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req UpdateDailyInstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updatedInstance, err := h.goalRepo.UpdateDailyInstance(instance.ID, userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedInstance)
}

// HandleGetGoalHistory godoc
// @Summary Get goal history
// @Description Get daily instances for a goal within a date range
// @Tags goals
// @Produce json
// @Security BearerAuth
// @Param goalId path string true "Goal ID"
// @Param startDate query string false "Start date (YYYY-MM-DD format, defaults to 30 days ago)"
// @Param endDate query string false "End date (YYYY-MM-DD format, defaults to today)"
// @Success 200 {array} DailyGoalInstance
// @Failure 400 {string} string "Invalid date format"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Goal not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/goals/{goalId}/history [get]
func (h *Handlers) HandleGetGoalHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	goalID := r.URL.Path[len("/api/goals/"):]
	goalID = goalID[:len(goalID)-len("/history")]
	if goalID == "" {
		http.Error(w, "Goal ID is required", http.StatusBadRequest)
		return
	}

	_, err := h.goalRepo.GetGoalByID(goalID, userID)
	if err != nil {
		if err.Error() == "goal not found" {
			http.Error(w, "Goal not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	endDate := time.Now()
	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		endDate = parsed
	}

	startDate := endDate.AddDate(0, 0, -30)
	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		startDate = parsed
	}

	instances, err := h.goalRepo.GetDailyInstancesByGoal(goalID, userID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instances)
}