package api

import (
	"net/http"

	"github.com/JoshPugli/grindhouse-api/internal/auth"
	"github.com/JoshPugli/grindhouse-api/internal/goals"
	"github.com/JoshPugli/grindhouse-api/internal/middleware"
)

// constructor is responsible for all the top-level HTTP stuff that applies to all endpoints,
// like CORS, auth middleware, and logging
func NewServer() http.Handler {
	mux := http.NewServeMux()

	authHandlers := auth.NewAuthHandlers()
	goalsStore := goals.NewStore()
	goalsHandlers := goals.NewHandlers(goalsStore)

	addRoutes(mux, authHandlers)
	addGoalRoutes(mux, goalsHandlers)

	return middleware.CORS(mux)
}
