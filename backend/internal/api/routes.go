// Package api provides the main API routing and handlers
package api

import (
	"fmt"
	"net/http"

	"github.com/JoshPugli/grindhouse-api/internal/auth"
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
