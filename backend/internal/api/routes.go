package api

import (
	"fmt"
	"net/http"

	"github.com/JoshPugli/grindhouse-api/internal/auth"
	"github.com/JoshPugli/grindhouse-api/internal/goals"
)

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromContext(r.Context())
	fmt.Fprintf(w, "Protected route accessed by user: %s", userID)
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
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

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
