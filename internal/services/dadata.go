package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"wifionAutolead/internal/models"
)

type Dadata struct {
	apiKey string
}

func NewDadata(apiKey string) *Dadata {
	return &Dadata{apiKey: apiKey}
}

func (d *Dadata) GetInfoOnAddress(address string) (models.AddressResponseDadata, error) {
	const url string = "http://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address"

	// Создание body запроса
	data := map[string]string{
		"query": address,
		"count": "1",
	}

	// Сериализация в JSON
	body, err := json.Marshal(data)
	if err != nil {
		return models.AddressResponseDadata{}, fmt.Errorf("ошибка при преобразовании в JSON: %v", err)
	}

	// Создание запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return models.AddressResponseDadata{}, fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	// Добавление заголовков
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", d.apiKey))

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.AddressResponseDadata{}, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Получение ответа
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return models.AddressResponseDadata{}, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	// Преобразование в нужный формат
	var response models.AddressResponseDadata
	err = json.Unmarshal(body, &response)
	if err != nil {
		return models.AddressResponseDadata{}, fmt.Errorf("ошибка при преобразовании в JSON: %v", err)
	}

	// Проверка на наличие данных
	if len(response.Suggestions) == 0 {
		return models.AddressResponseDadata{}, fmt.Errorf("нет предложений по адресу")
	}

	response.Suggestions[0].Data.RegionID = response.Suggestions[0].Data.RegionID[0:2]
	if response.Suggestions[0].Data.Flat == "" {
		response.Suggestions[0].Data.Flat = "0"
	}

	return response, nil
}
