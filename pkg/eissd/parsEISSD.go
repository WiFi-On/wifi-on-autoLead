package eissd

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	eissd "wifionAutolead/pkg/eissd/models"
)

// EISSDpars - структура для работы с EISSD
type EISSDpars struct {
	URL string
	sertPath string
	sertKey string
}
// NewEISSDpars создает новый экземпляр EISSDpars
func NewEISSDpars(url string, sertPath string, sertKey string) *EISSDpars {
	return &EISSDpars{URL: url, sertPath: sertPath, sertKey: sertKey}
}

// getClient создает и возвращает HTTP-клиент с настройками TLS
func createTLSClient(certPath, keyPath string) (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке сертификата: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &http.Client{
		Transport: transport,
	}, nil
}
// sendRequest создает и возвращает HTTP-запрос
func sendRequest(client *http.Client, url, requestBody string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Content-Type", "text/xml")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
// GetDistrictsOrAddresses получает список адресов или населенных пунктов по региону
func (e *EISSDpars) GetDistrictsOrAddresses(regionID string, structAddrObject int) ([]eissd.Address, error) {
	dateRequest := time.Now().UTC().Format("2006-01-02T15:04:05+00:00")

	requestBody := fmt.Sprintf(`
		<GetAddressInfoAgent DateRequest="%s" IdRequest="10001">
			<Release>2</Release>
			<RegionId>%s</RegionId>
			<StructAddrObject>%d</StructAddrObject>
		</GetAddressInfoAgent>`, dateRequest, regionID, structAddrObject)

	client, err := createTLSClient(e.sertPath, e.sertKey)
	if err != nil {
		return nil, err
	}

	body, err := sendRequest(client, "https://mpz.rt.ru/xmlInteface", requestBody)
	if err != nil {
		return nil, err
	}

	var result eissd.GetAddressInfoAgentResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге XML: %w", err)
	}

	return result.Addresses, nil
}
// GetHouses получает список домов по региону
func (e *EISSDpars) GetHouses(regionID string) ([]eissd.AddressHouse, error) {
	dateRequest := time.Now().UTC().Format("2006-01-02T15:04:05+00:00")

	requestBody := fmt.Sprintf(`
		<GetAddressHouseInfoAgent DateRequest="%s" IdRequest="10001">
			<Release>2</Release>
			<RegionId>%s</RegionId>
		</GetAddressHouseInfoAgent>`, dateRequest, regionID)

	client, err := createTLSClient(e.sertPath, e.sertKey)
	if err != nil {
		return nil, err
	}
	
	body, err := sendRequest(client, "https://mpz.rt.ru/xmlInteface", requestBody)
	if err != nil {
		return nil, err
	}

	// Парсинг XML в структуру
	var result eissd.GetAddressHouseInfoAgentResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге XML: %w", err)
	}

	return result.AddressHouses, nil
}
