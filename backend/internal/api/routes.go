// Package api provides the main API routing and handlers
package api

import (
	"fmt"
	"net/http"

	"github.com/JoshPugli/grindhouse-api/internal/auth"
	"github.com/JoshPugli/grindhouse-api/internal/goals"
	"github.com/JoshPugli/grindhouse-api/internal/user"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// protectedHandler godoc
// @Summary Protected endpoint
// @Description Example protected endpoint that requires authentication
// @Tags protected
// @Produce plain
// @Security BearerAuth
// @Success 200 {string} string "Protected route accessed"
// @Failure 401 {string} string "Unauthorized"
// @Router /api/protected [get]
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromContext(r.Context())
	fmt.Fprintf(w, "Protected route accessed by user: %s", userID)
}

// healthHandler godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /api/health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func addRoutes(
	mux *http.ServeMux,
	userHandlers *user.Handlers,
	goalHandlers *goals.Handlers,
) {
	// Public auth routes
	mux.HandleFunc("/api/auth/login", userHandlers.HandleLogin)
	mux.HandleFunc("/api/auth/register", userHandlers.HandleRegister)
	
	// Protected routes
	mux.Handle("/api/auth/me", auth.AuthMiddleware(http.HandlerFunc(userHandlers.HandleMe)))
	mux.Handle("/api/protected", auth.AuthMiddleware(http.HandlerFunc(protectedHandler)))
	
	// Goal routes
	mux.Handle("/api/goals", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			goalHandlers.HandleCreateGoal(w, r)
		case http.MethodGet:
			goalHandlers.HandleGetGoals(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	mux.Handle("/api/goals/today", auth.AuthMiddleware(http.HandlerFunc(goalHandlers.HandleGetGoalsToday)))
	mux.Handle("/api/goals/", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/api/goals/" {
			switch r.Method {
			case http.MethodGet:
				goalHandlers.HandleGetGoals(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else if len(path) > 12 && path[len(path)-8:] == "/history" {
			goalHandlers.HandleGetGoalHistory(w, r)
		} else if len(path) > 10 && path[len(path)-6:] == "/daily" {
			goalHandlers.HandleUpdateDailyInstance(w, r)
		} else {
			switch r.Method {
			case http.MethodGet:
				goalHandlers.HandleGetGoal(w, r)
			case http.MethodPut:
				goalHandlers.HandleUpdateGoal(w, r)
			case http.MethodDelete:
				goalHandlers.HandleDeleteGoal(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})))
	
	// Public routes
	mux.HandleFunc("/api/health", healthHandler)
	
	// Swagger documentation
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	
	mux.Handle("/", http.NotFoundHandler())
}
