package api

import (
	"fmt"
	"net/http"

	"github.com/JoshPugli/grindhouse-api/internal/auth"
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
