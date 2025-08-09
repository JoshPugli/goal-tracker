package api

import (
	"net/http"
)

// constructor is responsible for all the top-level HTTP stuff that applies to all endpoints, 
// like CORS, auth middleware, and logging
func NewServer() *http.ServeMux {
	mux := http.NewServeMux()
	
	addRoutes(mux)

	return mux
}