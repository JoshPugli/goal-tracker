package api

import (
	"log"
	"net/http"

	"github.com/JoshPugli/grindhouse-api/internal/database"
	"github.com/JoshPugli/grindhouse-api/internal/middleware"
	"github.com/JoshPugli/grindhouse-api/internal/user"
	
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

	userRepo := user.NewRepository(db)
	userHandlers := user.NewHandlers(userRepo)

	addRoutes(mux, userHandlers)

	return middleware.CORS(mux)
}
