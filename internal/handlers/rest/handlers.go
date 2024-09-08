package handlers

import (
	"encoding/json"
	"net/http"
	"wifionAutolead/internal/models"
	"wifionAutolead/internal/services"
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

func GetTHVHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Проверка метода запроса
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверный метод запроса"})
		return
	}

	// Проверка на пустое тело запроса
	if r.ContentLength == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Тело запроса пустое"})
		return
	}

	// Декодирование JSON тела запроса в структуру AddressData
	var address models.AddressData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&address); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверное тело запроса"})
		return
	}
	defer r.Body.Close()

	// Проверка, что значение адреса не пустое
	if address.Address == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Значение адреса пустое"})
		return
	}

	// Вызов функции бизнес-логики CheckTHV с адресом
	success, message, err := services.CheckTHV(address.Address)
	if err != nil {
		// В случае ошибки выполнения CheckTHV возвращаем ошибку
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	// Формируем ответ на основе результата CheckTHV
	response := models.CheckTHVResponse{
		Success: success,
		Message: message,
	}

	// Отправляем ответ клиенту
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Не удалось сформировать ответ"})
	}
}
