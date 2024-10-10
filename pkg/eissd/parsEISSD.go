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

type EISSDpars struct {
	sertPath string
	sertKey string
	urlCRM string
}

func NewEISSDpars(sertPath string, sertKey string, urlCRM string) *EISSDpars {
	return &EISSDpars{sertPath: sertPath, sertKey: sertKey, urlCRM: urlCRM}
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
// sendRequest отправляет запрос на сервер и возвращает полученный ответ
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
// GetDistrictsOrAddresses получает список дистриктов или улиц в указанном регионе
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

	body, err := sendRequest(client, e.urlCRM, requestBody)
	if err != nil {
		return nil, err
	}

	var result eissd.GetAddressInfoAgentResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге XML: %w", err)
	}

	return result.Addresses, nil
}
// GetTarrifsOnRegion получает список тарифов в указанном регионе
func (e *EISSDpars) GetTarrifsOnRegion(region string) ([]eissd.GetTariffPlansAgent, error) {

	requestBody := fmt.Sprintf(`
		<GetTariffPlansAgent>
    		<RegionId>%s</RegionId>
		</GetTariffPlansAgent>`, region)

	client, err := createTLSClient(e.sertPath, e.sertKey)
	if err != nil {
		return nil, err
	}
	
	body, err := sendRequest(client, e.urlCRM, requestBody)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	var data []eissd.GetTariffPlansAgent
	err = xml.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("ошибка при разборе XML: %w", err)
	}

	return data, nil
}
// GetTariffsMVNO получает список тарифов MVNO
func (e *EISSDpars) GetTariffsMVNO() ([]eissd.GetTariffPlansMvno, error) {
	requestBody := "<GetTariffPlansMvno></GetTariffPlansMvno>"
		
	client, err := createTLSClient(e.sertPath, e.sertKey)
	if err != nil {
		return nil, err
	}
	
	body, err := sendRequest(client, e.urlCRM, requestBody)
	if err != nil {
		return nil, err
	}

	var data []eissd.GetTariffPlansMvno
	err = xml.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("ошибка при разборе XML: %w", err)
	}
	
	return data, nil
}
