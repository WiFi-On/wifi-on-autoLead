package routes

import (
	"net/http"
	handlers "wifionAutolead/internal/controllers/rest"
)

// RegisterRoutes регистрирует маршруты и соответствующие хэндлеры
func RegisterRoutes(mux *http.ServeMux, eissdHandler *handlers.EISSDController) {
	mux.HandleFunc("/api/v1/rtkCRM/getTHVonAddress", eissdHandler.GetTHVHandler)
}
