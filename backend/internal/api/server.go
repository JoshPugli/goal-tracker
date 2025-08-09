package api

import (
	"log"
	"net/http"

	"github.com/JoshPugli/grindhouse-api/internal/auth"
	"github.com/JoshPugli/grindhouse-api/internal/database"
	"github.com/JoshPugli/grindhouse-api/internal/middleware"
	"github.com/JoshPugli/grindhouse-api/internal/repository"
	
	_ "github.com/JoshPugli/grindhouse-api/docs"
)

// constructor is responsible for all the top-level HTTP stuff that applies to all endpoints,
// like CORS, auth middleware, and logging
func NewServer() http.Handler {
	mux := http.NewServeMux()
	
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	authHandlers := auth.NewAuthHandlers(userRepo)

	addRoutes(mux, authHandlers)

	return middleware.CORS(mux)
}
