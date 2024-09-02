package handlers

import (
	"encoding/json"
	"net/http"
	"wifionAutolead/internal/models"
)

func TechnicalFeasibilityCheckHandler(w http.ResponseWriter, r *http.Request) {
    // Устанавливаем заголовок Content-Type
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверный метод запроса"})
        return
    }

    if r.ContentLength == 0 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Тело запроса пустое"})
        return
    }

    var address models.AddressData
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&address); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверное тело запроса"})
        return
    }
    defer r.Body.Close()

    if address.Address == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Значение адреса пустое"})
        return
    }

    response := models.AddressData{
        Address: address.Address,
    }

    if err := json.NewEncoder(w).Encode(response); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Не удалось сформировать ответ"})
        return
    }
}
func GetTHV(w http.ResponseWriter, r *http.Request) {


}