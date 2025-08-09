package api

import (
	"net/http"

	"github.com/JoshPugli/grindhouse-api/internal/auth"
	"github.com/JoshPugli/grindhouse-api/internal/middleware"
)

// constructor is responsible for all the top-level HTTP stuff that applies to all endpoints,
// like CORS, auth middleware, and logging
func NewServer() http.Handler {
	mux := http.NewServeMux()
	
	authHandlers := auth.NewAuthHandlers()

	addRoutes(mux, authHandlers)

	return middleware.CORS(mux)
}
