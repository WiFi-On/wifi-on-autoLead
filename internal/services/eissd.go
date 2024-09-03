package services

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"wifionAutolead/internal/models"
)

// GetClient возвращает клиент для запроса к EISSD
func getClient(cert tls.Certificate) *http.Client {
	// Создание конфигурации TLS
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	// Настройка транспортного уровня с использованием TLS
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &http.Client{
		Transport: transport,
	}
}

// FetchAddressDirectory извлекает справочник населенных пунктов или адресов
func FetchAddressDirectory(regionID int, structAddrObject int, searchName string) ([]models.Address, error) {
	// Форматирование текущего времени
	dateRequest := time.Now().UTC().Format("2006-01-02T15:04:05+00:00")

	// Создание XML тела запроса
	requestBody := fmt.Sprintf(`
		<GetAddressInfoAgent DateRequest="%s" IdRequest="10001">
			<Release>2</Release>
			<RegionId>%d</RegionId>
			<StructAddrObject>%d</StructAddrObject>
		</GetAddressInfoAgent>`, dateRequest, regionID, structAddrObject)

	// Загрузка сертификата и ключа
	cert, err := tls.LoadX509KeyPair("../../common/certs/krivoshein.crt", "../../common/certs/krivoshein.key")
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке сертификата: %w", err)
	}

	client := getClient(cert)

	// Создание HTTP запроса
	req, err := http.NewRequest("POST", "https://mpz.rt.ru/xmlInteface", strings.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Установка необходимых заголовков
	req.Header.Set("Content-Type", "text/xml")

	// Отправка HTTP запроса
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении тела ответа: %w", err)
	}

	// Парсинг XML в структуру
	var result models.GetAddressInfoAgentResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге XML: %w", err)
	}

	// Фильтрация результатов по NameAddrObject
	var filteredAddresses []models.Address
	for _, address := range result.Addresses {
		if address.NameAddrObject == searchName {
			filteredAddresses = append(filteredAddresses, address)
		}
	}

	return filteredAddresses, nil
}

// FetchAddressHouseInfo ищет дома по указанному региону, StreetId и House
func FetchAddressHouseInfo(regionID int, searchStreetId, searchHouse string) ([]models.AddressHouse, error) {
	// Форматирование текущего времени
	dateRequest := time.Now().UTC().Format("2006-01-02T15:04:05+00:00")

	// Создание XML тела запроса
	requestBody := fmt.Sprintf(`
		<GetAddressHouseInfoAgent DateRequest="%s" IdRequest="10001">
			<Release>2</Release>
			<RegionId>%d</RegionId>
		</GetAddressHouseInfoAgent>`, dateRequest, regionID)

	// Загрузка сертификата и ключа
	cert, err := tls.LoadX509KeyPair("../../common/certs/krivoshein.crt", "../../common/certs/krivoshein.key")
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке сертификата: %w", err)
	}

	client := getClient(cert)

	// Создание HTTP запроса
	req, err := http.NewRequest("POST", "https://mpz.rt.ru/xmlInteface", strings.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Установка необходимых заголовков
	req.Header.Set("Content-Type", "text/xml")

	// Отправка HTTP запроса
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении тела ответа: %w", err)
	}

	// Парсинг XML в структуру
	var result models.GetAddressHouseInfoAgentResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге XML: %w", err)
	}

	// Фильтрация результатов по StreetId и House
	var filteredHouses []models.AddressHouse
	for _, house := range result.AddressHouses {
		if house.StreetId == searchStreetId && house.House == searchHouse {
			filteredHouses = append(filteredHouses, house)
		}
	}

	return filteredHouses, nil
}

// CheckConnectionPossibilityAgent выполняет проверку возможности подключения
func CheckConnectionPossibilityAgent(regionID int, cityID string, streetID string, houseID string, svcClassId int) (int, string, error) {
	// Форматирование текущего времени
	dateRequest := time.Now().UTC().Format("2006-01-02T15:04:05+00:00")

	requestBody := fmt.Sprintf(`
		<CheckConnectionPossibilityAgent DateRequest="%s" IdRequest="10001">
			<Release>2</Release>
			<RegionId>%d</RegionId>
			<CityId>%s</CityId>
			<StreetId>%s</StreetId>
			<HouseId>%s</HouseId>
			<SvcClassIds>
				<SvcClassId>%d</SvcClassId>
			</SvcClassIds>
		</CheckConnectionPossibilityAgent>`, dateRequest, regionID, cityID, streetID, houseID, svcClassId)

	// Загрузка сертификата и ключа
	cert, err := tls.LoadX509KeyPair("../../common/certs/krivoshein.crt", "../../common/certs/krivoshein.key")
	if err != nil {
		return 0, "", fmt.Errorf("ошибка при загрузке сертификата: %w", err)
	}

	client := getClient(cert)

	// Создание HTTP-запроса
	req, err := http.NewRequest("POST", "https://mpz.rt.ru/xmlInteface", strings.NewReader(requestBody))
	if err != nil {
		return 0, "", fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Установка необходимых заголовков
	req.Header.Set("Content-Type", "text/xml")

	// Отправка HTTP-запроса
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("ошибка при чтении тела ответа: %w", err)
	}

	// Парсинг XML в структуру
	var result models.CheckConnectionPossibilityResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return 0, "", fmt.Errorf("ошибка при парсинге XML: %w", err)
	}

	// Возврат результата
	return result.Response, result.Message, nil
}

// GetTarrifsOnRegion получение тарифов по региону
func GetTarrifsOnRegion(region int) (models.GetTariffPlansAgent, error) {

	requestBody := fmt.Sprintf(`
		<GetTariffPlansAgent>
    		<RegionId>%d</RegionId>
		</GetTariffPlansAgent>`, region)

	// Загрузка сертификата и ключа
	cert, err := tls.LoadX509KeyPair("../../common/certs/krivoshein.crt", "../../common/certs/krivoshein.key")
	if err != nil {
		return models.GetTariffPlansAgent{}, fmt.Errorf("ошибка при загрузке сертификата: %w", err)
	}

	client := getClient(cert)

	// Создание HTTP-запроса
	req, err := http.NewRequest("POST", "https://mpz.rt.ru/xmlInteface", strings.NewReader(requestBody))
	if err != nil {
		return models.GetTariffPlansAgent{}, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Установка необходимых заголовков
	req.Header.Set("Content-Type", "text/xml")

	// Отправка HTTP-запроса
	resp, err := client.Do(req)
	if err != nil {
		return models.GetTariffPlansAgent{}, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.GetTariffPlansAgent{}, fmt.Errorf("ошибка при чтении тела ответа: %w", err)
	}

	// Парсинг XML в структуру
	var data models.GetTariffPlansAgent
	err = xml.Unmarshal(body, &data)
	if err != nil {
		return models.GetTariffPlansAgent{}, fmt.Errorf("ошибка при разборе XML: %w", err)
	}

	// Возврат результата
	return data, nil
}
