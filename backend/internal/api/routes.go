// Package api provides the main API routing and handlers
package api

import (
	"fmt"
	"net/http"

	"github.com/JoshPugli/grindhouse-api/internal/auth"
	"github.com/JoshPugli/grindhouse-api/internal/goals"
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
	authHandlers *auth.AuthHandlers,
) {
	// Public auth routes
	mux.HandleFunc("/api/auth/login", authHandlers.HandleLogin)
	mux.HandleFunc("/api/auth/register", authHandlers.HandleRegister)

	// Protected routes
	mux.Handle("/api/auth/me", auth.AuthMiddleware(http.HandlerFunc(authHandlers.HandleMe)))
	mux.Handle("/api/protected", auth.AuthMiddleware(http.HandlerFunc(protectedHandler)))

	// Public routes
	mux.HandleFunc("/api/health", healthHandler)
	
	// Swagger documentation
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.Handle("/", http.NotFoundHandler())
}

func addGoalRoutes(
	mux *http.ServeMux,
	goalHandlers *goals.Handlers,
) {
	// Goals catalog
	mux.HandleFunc("/api/goals", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		goalHandlers.HandleListGoals(w, r)
	})

	// Today state
	mux.HandleFunc("/api/goals/today", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		goalHandlers.HandleToday(w, r)
	})

	// Stats by window
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		goalHandlers.HandleStats(w, r)
	})

	// Dashboard aggregate
	mux.HandleFunc("/api/dashboard", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		goalHandlers.HandleDashboard(w, r)
	})

	// Toggle complete for today
	// POST /api/goals/{id}/complete, DELETE /api/goals/{id}/complete
	mux.HandleFunc("/api/goals/", func(w http.ResponseWriter, r *http.Request) {
		// minimal pattern match for "/api/goals/{id}/complete"
		if len(r.URL.Path) >= len("/api/goals/") && hasSuffix(r.URL.Path, "/complete") {
			goalHandlers.HandleToggleComplete(w, r)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
}

// local helper to avoid new dep
func hasSuffix(s, suf string) bool {
	if len(suf) > len(s) {
		return false
	}
	return s[len(s)-len(suf):] == suf
}
