package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Data struct {
	Value string `json:"value"`
}

type Response struct {
	Suggestions []Data `json:"suggestions"`
}

type ResponseTHV struct {
	Suggestions []struct {
		Data struct {
			RegionKladrID string `json:"region_kladr_id"`
			Settlement    string `json:"settlement"`
			City          string `json:"city"`
			Street        string `json:"street"`
			House         string `json:"house"`
		} `json:"data"`
	} `json:"suggestions"`
}

func GetInfoOnAddress(address string, apiKey string) ([]byte, error) {
	const url string = "http://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address"

	// Создание body запроса
	data := map[string]string{
		"query": address,
		"count": "1",
	}

	// Сериализация в JSON
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("ошибка при преобразовании в JSON: %v", err)
	}

	// Создание запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	// Добавление заголовков
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apiKey))

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Получение ответа
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при преобразовании в JSON: %v", err)
	}

	// Преобразование в нужный формат
	var response Response
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, fmt.Errorf("ошибка при преобразовании в JSON: %v", err)
	}

	return []byte(response.Suggestions[0].Value), nil
}

func GetInfoOnAddressTHV(address string, apiKey string) (string, error) {
	const url string = "http://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address"

	// Создание body запроса
	data := map[string]string{
		"query": address,
		"count": "1",
	}

	// Сериализация в JSON
	body, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("ошибка при преобразовании в JSON: %v", err)
	}

	// Создание запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	// Добавление заголовков
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apiKey))

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Получение ответа
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	// Преобразование в нужный формат
	var response ResponseTHV
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("ошибка при десериализации JSON: %v", err)
	}

	if len(response.Suggestions) == 0 {
		return "", fmt.Errorf("нет предложений по данному адресу")
	}

	// Извлечение нужных данных
	dataInfo := response.Suggestions[0].Data

	// Первые две цифры region_kladr_id
	regionID := ""
	if len(dataInfo.RegionKladrID) >= 2 {
		regionID = dataInfo.RegionKladrID[:2]
	}

	// Логика для выбора между settlement и city
	settlementOrCity := dataInfo.Settlement
	if settlementOrCity == "" {
		settlementOrCity = dataInfo.City
	}

	// Улица и дом
	street := dataInfo.Street
	house := dataInfo.House

	// Формирование результирующей строки
	result := fmt.Sprintf("%s, %s, %s, %s", regionID, settlementOrCity, street, house)

	return result, nil
}
