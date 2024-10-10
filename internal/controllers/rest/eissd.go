package controllers

import (
	"encoding/json"
	"net/http"
	"wifionAutolead/internal/models"
	"wifionAutolead/internal/services"
)

type EISSDController struct {
	eissdService *services.EISSD
	bitrixService *services.Bitrix
	dadataService *services.Dadata
}

func NewEISSDController(eissdService *services.EISSD, bitrixService *services.Bitrix, dadataService *services.Dadata) *EISSDController {
	return &EISSDController{eissdService: eissdService, bitrixService: bitrixService, dadataService: dadataService}
}

func (h *EISSDController) GetTHVHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверный метод запроса"})
		return
	}

	queryParams := r.URL.Query()

	address := queryParams.Get("address")
	
	if len(queryParams) != 1 || address == ""  {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверные параметры запроса"})
		return
	}

	thv, districtFiasId, err := h.eissdService.CheckTHV(address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Ошибка в получении данных о ТХВ"})
		return
	}

	result := struct {
		Thv            []models.ReturnDataConnectionPos
		DistrictFiasId string
	}{
		Thv:            thv,
		DistrictFiasId: districtFiasId,
	}

	if result.Thv != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Ошибка в получении данных о ТХВ"})
		return
	}
}
