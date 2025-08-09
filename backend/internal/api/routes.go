package api

import (
	"net/http"
	"fmt"
)


func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request sent to: %s", r.URL.Path)
}
	

func addRoutes(
	mux                 *http.ServeMux,
	// logger              *logging.Logger,
) {
	// mux.Handle("/api/v1/", handleTenantsGet(logger, tenantsStore))
	// mux.Handle("/oauth2/", handleOAuth2Proxy(logger, authProxy))
	// mux.HandleFunc("/healthz", handleHealthzPlease(logger))
	mux.HandleFunc("/handle/", handler)
	mux.Handle("/", http.NotFoundHandler())
}