package routes

import (
	"net/http"
	handlers "wifionAutolead/internal/handlers/rest"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/technicalFeasibilityCheck", handlers.TechnicalFeasibilityCheckHandler)
	mux.HandleFunc("/getTHV", handlers.GetTHVHandler)
}
