package services

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Address содержит информацию об одном адресе из XML-ответа
type Address struct {
	RegionId       int    `xml:"RegionId"`
	AddrObjectId   string `xml:"AddrObjectId"`
	NameAddrObject string `xml:"NameAddrObject"`
	AbbrNameObject string `xml:"AbbrNameObject"`
	ParentId       string `xml:"ParentId"`
}

// GetAddressInfoAgentResponse структура для обработки всего XML-ответа
type GetAddressInfoAgentResponse struct {
	Addresses []Address `xml:"addresss>address"`
}

// FetchAddressDirectory извлекает справочник населенных пунктов или адресов
func FetchAddressDirectory(regionID int, structAddrObject int, searchName string) ([]Address, error) {
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

	// Создание конфигурации TLS
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	// Настройка транспортного уровня с использованием TLS
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: transport,
	}

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
	var result GetAddressInfoAgentResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге XML: %w", err)
	}

	// Фильтрация результатов по NameAddrObject
	var filteredAddresses []Address
	for _, address := range result.Addresses {
		if address.NameAddrObject == searchName {
			filteredAddresses = append(filteredAddresses, address)
		}
	}

	return filteredAddresses, nil
}

// AddressHouse содержит информацию об одном доме из XML-ответа
type AddressHouse struct {
	RegionId int    `xml:"RegionId"`
	StreetId string `xml:"StreetId"`
	HouseId  string `xml:"HouseId"`
	House    string `xml:"House"`
}

// GetAddressHouseInfoAgentResponse структура для обработки всего XML-ответа
type GetAddressHouseInfoAgentResponse struct {
	AddressHouses []AddressHouse `xml:"AddressHouses>AddressHouse"`
}

// FetchAddressHouseInfo ищет дома по указанному региону, StreetId и House
func FetchAddressHouseInfo(regionID int, searchStreetId, searchHouse string) ([]AddressHouse, error) {
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

	// Создание конфигурации TLS
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	// Настройка транспортного уровня с использованием TLS
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: transport,
	}

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
	var result GetAddressHouseInfoAgentResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге XML: %w", err)
	}

	// Фильтрация результатов по StreetId и House
	var filteredHouses []AddressHouse
	for _, house := range result.AddressHouses {
		if house.StreetId == searchStreetId && house.House == searchHouse {
			filteredHouses = append(filteredHouses, house)
		}
	}

	return filteredHouses, nil
}

// CheckConnectionPossibilityResponse описывает ответ на запрос проверки возможности подключения
type CheckConnectionPossibilityResponse struct {
	Response int    `xml:"Response"`
	Message  string `xml:"Message"`
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

	// Создание конфигурации TLS
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	// Настройка транспортного уровня с использованием TLS
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: transport,
	}

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
	var result CheckConnectionPossibilityResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return 0, "", fmt.Errorf("ошибка при парсинге XML: %w", err)
	}

	// Возврат результата
	return result.Response, result.Message, nil
}
