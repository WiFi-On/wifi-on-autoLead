package services

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"wifionAutolead/internal/models"
	"wifionAutolead/pkg/eissd"
	"wifionAutolead/pkg/utils"
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

// CheckConnectionPossibilityAgent выполняет проверку возможности подключения
func CheckConnectionPossibilityAgent(regionID int, cityID int, streetID int, houseID int, svcClassId int) (int, string, error) {
	// Форматирование текущего времени
	dateRequest := time.Now().UTC().Format("2006-01-02T15:04:05+00:00")

	requestBody := fmt.Sprintf(`
		<CheckConnectionPossibilityAgent DateRequest="%s" IdRequest="10001">
			<Release>2</Release>
			<RegionId>%d</RegionId>
			<CityId>%d</CityId>
			<StreetId>%d</StreetId>
			<HouseId>%d</HouseId>
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

// CheckTHV проверяет возможность подключения и возвращает успешное подключение или сообщение об ошибке.
func CheckTHV(address string) (bool, string, error) {
	// Загрузка .env файла и получение API ключа
	err := godotenv.Load("../common/conf/.env")
	if err != nil {
		return false, "", fmt.Errorf("ошибка загрузки .env файла: %v", err)
	}

	apiKey := os.Getenv("DADATA_API_KEY")
	if apiKey == "" {
		return false, "", fmt.Errorf("API ключ не найден в .env файле")
	}

	// Вызов внешнего сервиса для получения информации по адресу
	resultDaData, err := utils.GetInfoOnAddressTHV(address, apiKey)
	if err != nil {
		return false, "", fmt.Errorf("ошибка при получении данных по адресу: %v", err)
	}

	// Разделяем строку адреса по запятым
	parts := strings.Split(resultDaData, ",")
	if len(parts) != 4 {
		return false, "", fmt.Errorf("неверный формат адреса. Ожидается формат: 'regionID, cityName, streetName, houseNumber'")
	}

	// Присваиваем переменным значения, полученные из строки адреса
	regionID := parts[0]
	cityName := strings.TrimSpace(parts[1])
	streetName := strings.TrimSpace(parts[2])
	houseNumber := strings.TrimSpace(parts[3])

	// Преобразование regionID из string в int
	regionIDInt, err := strconv.Atoi(regionID)
	if err != nil {
		return false, "", fmt.Errorf("ошибка преобразования regionID: %v", err)
	}

	// Преобразование houseNumber из string в int
	houseNumberInt, err := strconv.Atoi(houseNumber)
	if err != nil {
		return false, "", fmt.Errorf("ошибка преобразования houseNumber: %v", err)
	}

	// Подключение к БД
	eissdDB, err := eissd.NewDB("../common/db/eissd.db")
	if err != nil {
		return false, "", fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	// Получение cityID (или districtID) из БД
	cityID, err := eissdDB.GetDistrictIDByRegionAndName(regionIDInt, cityName)
	if err != nil {
		return false, "", fmt.Errorf("ошибка получения districtID по региону и имени города: %v", err)
	}

	// Получение streetID по региону, cityID и имени улицы
	streetID, err := eissdDB.GetStreetIDByRegionNameAndDistrict(cityID, streetName, regionIDInt)
	if err != nil {
		return false, "", fmt.Errorf("ошибка получения streetID по региону и имени улицы: %v", err)
	}

	// Получение houseID по региону, streetID и номеру дома
	houseID, err := eissdDB.GetHouseIDByRegionStreetAndHouse(regionIDInt, streetID, houseNumberInt)
	if err != nil {
		return false, "", fmt.Errorf("ошибка получения houseID по региону, улице и номеру дома: %v", err)
	}

	// Проверка возможности подключения
	responseCode, message, err := CheckConnectionPossibilityAgent(regionIDInt, cityID, streetID, houseID, 2)
	if err != nil {
		return false, "", fmt.Errorf("ошибка проверки возможности подключения: %v", err)
	}

	// Возвращаем результат на основе responseCode
	if responseCode == 0 {
		// Подключение возможно
		return true, "Подключение возможно", nil
	}

	// Подключение невозможно
	return false, message, nil
}
