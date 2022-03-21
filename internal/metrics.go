package internal

import (
	"net/http"

	"github.com/gorilla/mux"
)

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func AddMetricsRoute(router *mux.Router) {
	router.HandleFunc("/metrics", handleMetrics)
}
