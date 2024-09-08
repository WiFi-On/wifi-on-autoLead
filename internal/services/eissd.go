package services

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"wifionAutolead/internal/models"
	"wifionAutolead/pkg/eissd"
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

func CheckTHV(address string) {
	// Разделяем строку адреса по запятым
	parts := strings.Split(address, ",")
	if len(parts) != 4 {
		panic("Неверный формат адреса. Ожидается формат: 'regionID, cityName, streetName, houseNumber'")
	}

	// Присваиваем переменным значения, полученные из строки адреса
	regionID := parts[0]
	cityName := strings.TrimSpace(parts[1])
	streetName := strings.TrimSpace(parts[2])
	houseNumber := strings.TrimSpace(parts[3])

	// Преобразование regionID из string в int
	regionIDInt, err := strconv.Atoi(regionID)
	if err != nil {
		panic(fmt.Sprintf("Ошибка преобразования regionID: %v", err))
	}

	// Преобразование regionID из string в int
	houseNubmerInt, err := strconv.Atoi(houseNumber)
	if err != nil {
		panic(fmt.Sprintf("Ошибка преобразования regionID: %v", err))
	}

	eissdDB, err := eissd.NewDB("../common/db/eissd.db")
	if err != nil {
		panic(err)
	}

	// Вызов функции для получения CityID (или DistrictID) из БД
	cityID, err := eissdDB.GetDistrictIDByRegionAndName(regionIDInt, cityName)
	if err != nil {
		panic(err)
	}

	// Вызов функции для получения streetID по region, cityID и имени улицы
	streetID, err := eissdDB.GetStreetIDByRegionNameAndDistrict(cityID, streetName, regionIDInt)
	if err != nil {
		panic(err)
	}

	// Вызов функции для получения houseID по region, streetID и номеру дома
	houseID, err := eissdDB.GetHouseIDByRegionStreetAndHouse(regionIDInt, streetID, houseNubmerInt)
	if err != nil {
		panic(err)
	}

	// Вызов функции для проверки подключения
	responseCode, message, err := CheckConnectionPossibilityAgent(regionIDInt, cityID, streetID, houseID, 2)
	if err != nil {
		panic(err)
	}

	// Обработка результата
	if responseCode == 0 {
		// Подключение возможно
		fmt.Println("Подключение есть")
	} else {
		// Возвращаем сообщение об ошибке
		fmt.Printf("Ошибка: %s\n", message)
	}
}
