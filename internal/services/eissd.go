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

// FetchAddressDirectory извлекает справочник населенных пунктов или адресов и сохраняет в мапу
func FetchAddressDirectory(regionID int, structAddrObject int, searchName string) (map[string][]string, error) {
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

	// Создание мапы для хранения отфильтрованных данных
	addressMap := make(map[string][]string)

	// Использование XML Decoder для посимвольного чтения XML
	decoder := xml.NewDecoder(strings.NewReader(string(body)))
	var currentElement string
	var regionIDStr, nameAddrObject, addrObjectId, parentId string

	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("ошибка при чтении токенов XML: %w", err)
		}

		switch se := t.(type) {
		case xml.StartElement:
			currentElement = se.Name.Local
		case xml.CharData:
			data := strings.TrimSpace(string(se))
			if data == "" {
				continue
			}
			switch currentElement {
			case "RegionId":
				regionIDStr = data
			case "NameAddrObject":
				nameAddrObject = data
			case "AddrObjectId":
				addrObjectId = data
			case "ParentId":
				parentId = data
			}
		case xml.EndElement:
			if se.Name.Local == "address" {
				if nameAddrObject == searchName && regionIDStr == fmt.Sprintf("%d", regionID) {
					key := fmt.Sprintf("%s;%s", regionIDStr, nameAddrObject)
					value := fmt.Sprintf("%s;%s", addrObjectId, parentId)
					addressMap[key] = append(addressMap[key], value)
				}
				// Сброс значений для следующего Address
				regionIDStr, nameAddrObject, addrObjectId, parentId = "", "", "", ""
			}
		}
	}
	return addressMap, nil
}

func FetchAddressHouseInfo(regionID int, searchStreetId, searchHouse string) (map[string][]string, error) {
	// Замер времени начала выполнения функции
	startFunction := time.Now()

	// Форматирование текущего времени
	dateRequest := time.Now().UTC().Format("2006-01-02T15:04:05+00:00")

	// Создание XML тела запроса
	requestBody := fmt.Sprintf(`
		<GetAddressHouseInfoAgent DateRequest="%s" IdRequest="10001">
			<Release>2</Release>
			<RegionId>%d</RegionId>
		</GetAddressHouseInfoAgent>`, dateRequest, regionID)

	// Замер времени для загрузки сертификата
	startCert := time.Now()
	cert, err := tls.LoadX509KeyPair("../../common/certs/krivoshein.crt", "../../common/certs/krivoshein.key")
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке сертификата: %w", err)
	}
	fmt.Printf("Время загрузки сертификата: %v\n", time.Since(startCert))

	client := getClient(cert)

	// Создание HTTP запроса
	startRequest := time.Now()
	req, err := http.NewRequest("POST", "https://mpz.rt.ru/xmlInteface", strings.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	fmt.Printf("Время создания HTTP запроса: %v\n", time.Since(startRequest))

	// Установка необходимых заголовков
	req.Header.Set("Content-Type", "text/xml")

	// Отправка HTTP запроса
	startHTTP := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	fmt.Printf("Время отправки HTTP запроса и получения ответа: %v\n", time.Since(startHTTP))
	defer resp.Body.Close()

	// Чтение ответа
	startRead := time.Now()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении тела ответа: %w", err)
	}
	fmt.Printf("Время чтения ответа: %v\n", time.Since(startRead))

	// Создание мапы для хранения отфильтрованных данных
	houseMap := make(map[string][]string)

	// Использование XML Decoder для посимвольного чтения XML
	startParse := time.Now()
	decoder := xml.NewDecoder(strings.NewReader(string(body)))
	var currentElement string
	var regionIDStr, streetID, houseID, house string

	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("ошибка при чтении токенов XML: %w", err)
		}

		switch se := t.(type) {
		case xml.StartElement:
			currentElement = se.Name.Local
		case xml.CharData:
			data := strings.TrimSpace(string(se))
			if data == "" {
				continue
			}
			switch currentElement {
			case "RegionId":
				regionIDStr = data
			case "StreetId":
				streetID = data
			case "HouseId":
				houseID = data
			case "House":
				house = data
			}
		case xml.EndElement:
			if se.Name.Local == "AddressHouse" {
				if streetID == searchStreetId && house == searchHouse {
					key := fmt.Sprintf("%s;%s", regionIDStr, streetID)
					value := fmt.Sprintf("%s;%s", houseID, house)
					houseMap[key] = append(houseMap[key], value)
				}
				// Сброс значений для следующего AddressHouse
				regionIDStr, streetID, houseID, house = "", "", "", ""
			}
		}
	}
	fmt.Printf("Время парсинга XML: %v\n", time.Since(startParse))

	// Замер общего времени выполнения функции
	fmt.Printf("Общее время выполнения функции FetchAddressHouseInfo: %v\n", time.Since(startFunction))

	return houseMap, nil
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
