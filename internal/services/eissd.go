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
	"wifionAutolead/internal/repository"
)

type EISSD struct {
	repo      *repository.DB
	dadata    *Dadata
	logger    *LoggerService
	pathCert  string
	pathKey   string
	urlCRM    string
}

func NewEISSD(repo *repository.DB, pathCert string, pathKey string, urlCRM string, dadata *Dadata, logger *LoggerService) *EISSD {
	return &EISSD{repo: repo, pathCert: pathCert, pathKey: pathKey, urlCRM: urlCRM, dadata: dadata, logger: logger}
}
func getClient(cert tls.Certificate) *http.Client {
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &http.Client{
		Transport: transport,
	}
}
func (e *EISSD) CheckTHV(address string) ([]models.ReturnDataConnectionPos, string, error) {
	dateRequest := time.Now().UTC().Format("2006-01-02T15:04:05+00:00")
	dadataInfo, err := e.dadata.GetInfoOnAddress(address)
	if err != nil {
		e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
		e.logger.Log("ErrorAddress: " + address)
		return nil, "", fmt.Errorf("дадата невернула данные о адресе: %v", err)
	}
	infoAddressDadata := dadataInfo.Suggestions[0].Data
	
	//Блок условий для нахождения населенного пункта(district) 
	var idDistrict string
	var districtFiasID string
	if infoAddressDadata.Area != "" && infoAddressDadata.City != ""{
		areaId, err := e.repo.GetDistrictIDByRegionAndName(infoAddressDadata.RegionID, infoAddressDadata.Area)
		if err != nil {
			e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
			e.logger.Log("ErrorAddress: " + address)
			return nil, "", fmt.Errorf("не удалось получить айди района в регионе: %v", err)
		}

		cityId, err := e.repo.GetDistrictIDByParentIDandName(areaId, infoAddressDadata.City)
		if err != nil {
			e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
			e.logger.Log("ErrorAddress: " + address)
			return nil, "",fmt.Errorf("не удалось получить айди района: %v", err)
		}
		


		if cityId == "" {
			e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, "не удалось получить айди населенного пункта", address))
			e.logger.Log("ErrorAddress: " + address)
			return nil, "", fmt.Errorf("не удалось получить айди населенного пункта: %v", err)
		}
		idDistrict = cityId
		districtFiasID = infoAddressDadata.CityFiasId

	} else if infoAddressDadata.City != "" && infoAddressDadata.Settlement != ""{
		cityId, err := e.repo.GetDistrictIDByRegionAndName(infoAddressDadata.RegionID, infoAddressDadata.City)
		if err != nil {
			e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
			e.logger.Log("ErrorAddress: " + address)
			return nil, "", fmt.Errorf("не удалось получить айди города в регионе: %v", err)
		}
		
		settlementId, err := e.repo.GetDistrictIDByParentIDandName(cityId, infoAddressDadata.Settlement)
		if err != nil {
			e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
			e.logger.Log("ErrorAddress: " + address)
			return nil, "", fmt.Errorf("не удалось получить айди населенного пункта: %v", err)
		}



		if settlementId == "" {
			e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, "не удалось получить айди населенного пункта", address))
			e.logger.Log("ErrorAddress: " + address)
			return nil, "", fmt.Errorf("не удалось получить айди населенного пункта: %v", err)
		}
		idDistrict = settlementId
		districtFiasID = infoAddressDadata.SettlementFiasId

	} else if infoAddressDadata.Area  != "" && infoAddressDadata.Settlement != ""{
		areaId, err := e.repo.GetDistrictIDByRegionAndName(infoAddressDadata.RegionID, infoAddressDadata.Area)
		if err != nil {
			e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
			e.logger.Log("ErrorAddress: " + address)
			return nil, "", fmt.Errorf("не удалось получить айди района в регионе: %v", err)
		}
		
		settlementId, err := e.repo.GetDistrictIDByParentIDandName(areaId, infoAddressDadata.Settlement)
		if err != nil {
			e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
			e.logger.Log("ErrorAddress: " + address)
			return nil, "", fmt.Errorf("не удалось получить айди населенного пункта: %v", err)
		}



		if settlementId == "" {
			e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, "не удалось получить айди населенного пункта", address))
			e.logger.Log("ErrorAddress: " + address)
			return nil, "", fmt.Errorf("не удалось получить айди населенного пункта: %v", err)
		}
		idDistrict = settlementId
		districtFiasID = infoAddressDadata.SettlementFiasId

	} else if infoAddressDadata.City != "" {
		cityId, err := e.repo.GetDistrictIDByRegionAndName(infoAddressDadata.RegionID, infoAddressDadata.City)
		if err != nil {
			e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
			e.logger.Log("ErrorAddress: " + address)
			return nil, "", fmt.Errorf("не удалось получить айди города в регионе: %v", err)
		}



		if cityId == "" {
			e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, "не удалось получить айди населенного пункта", address))
			e.logger.Log("ErrorAddress: " + address)
			return nil, "", fmt.Errorf("не удалось получить айди населенного пункта: %v", err)
		}
		idDistrict = cityId
		districtFiasID = infoAddressDadata.CityFiasId

	} else {
		e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, "ошибка в нахождении district", address))
		e.logger.Log("ErrorAddress: " + address)
		return nil, "", fmt.Errorf("ошибка в нахождении district")
	}


	//Получение айди улицы 
	idStreet, err := e.repo.GetStreetIDByNameAndDistrictId(infoAddressDadata.Street, idDistrict)
	if err != nil  {
		e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
		e.logger.Log("ErrorAddress: " + address)
		return nil, "", fmt.Errorf("не удалось получить айди улицы: %v", err)
	}
	if idStreet == "" {
		e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, "ошибка в нахождении street", address))
		e.logger.Log("ErrorAddress: " + address)
		return nil, "", fmt.Errorf("ошибка в нахождении street")
	}
	//Получение айди дома
	houseAndBlock := infoAddressDadata.House + infoAddressDadata.Block
	idHouse, err := e.repo.GetHouseIDByStreetIdAndHouse(idStreet, houseAndBlock)
	if err != nil {
		e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
		e.logger.Log("ErrorAddress: " + address)
		return nil, "", fmt.Errorf("не удалось получить айди дома: %v", err)
	}
	if idHouse == "" {
		e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, "ошибка в нахождении house", address))
		e.logger.Log("ErrorAddress: " + address)
		return nil, "", fmt.Errorf("ошибка в нахождении house")
	}
	// Отправка запроса
	requestBody := fmt.Sprintf(`
		<CheckConnectionPossibilityAgent DateRequest="%s" IdRequest="10001">
    		<Release>2</Release>
    		<RegionId>%s</RegionId>
    		<CityId>%s</CityId>
    		<StreetId>%s</StreetId>
    		<HouseId>%s</HouseId>
    		<Flat>%s</Flat>
    		<TypeAdrId>0</TypeAdrId >
    		<SvcClassIds>
        		<SvcClassId>2</SvcClassId>
    		</SvcClassIds>
		</CheckConnectionPossibilityAgent>`, dateRequest, infoAddressDadata.RegionID, idDistrict, idStreet, idHouse, infoAddressDadata.Flat)
	cert, err := tls.LoadX509KeyPair(e.pathCert, e.pathKey)
	if err != nil {
		e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
		e.logger.Log("ErrorAddress: " + address)
		return nil, "", fmt.Errorf("ошибка при загрузке сертификата: %w", err)
	}
	client := getClient(cert)
	req, err := http.NewRequest("POST", e.urlCRM, strings.NewReader(requestBody))
	if err != nil {
		e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
		e.logger.Log("ErrorAddress: " + address)
		return nil, "", fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
		e.logger.Log("ErrorAddress: " + address)
		return nil, "", fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()
	// Чтение и разбор тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
		e.logger.Log("ErrorAddress: " + address)
		return nil, "", fmt.Errorf("ошибка при чтении тела ответа: %w", err)
	}
	var data models.CheckConnectionPossibilityAgent
	err = xml.Unmarshal(body, &data)
	if err != nil {
		e.logger.Log(fmt.Sprintf("time:%s || error:%s || address:%s", dateRequest, err.Error(), address))
		e.logger.Log("ErrorAddress: " + address)
		return nil, "", fmt.Errorf("ошибка при разборе XML: %w", err)
	}

	var thv []models.ReturnDataConnectionPos
	for _, item := range data.ConnectionPoss {
		thv = append(thv, models.ReturnDataConnectionPos{
			TechName: item.TechName,
			Res:      item.Res,
		})
	}

	// Запись в логи и возвращение результата
	e.logger.Log(fmt.Sprintf("time:%s || result:%s || address:%s", dateRequest, body, address))
	e.logger.Log("FoundAddress: " + address)
	return thv, districtFiasID, nil
}
func (e *EISSD) SendLead(idLead string, name string, surname string, patronymic string, phone string, xmlProduct string,
	idRegion string, idDistrict string, idStreet string, idHouse string, flat string) ([]models.ResponseSendLead, error) {

	dateRequest := time.Now().UTC().Format("2006-01-02T15:04:05+00:00")
	requestBody := fmt.Sprintf(`
		<CreatePacketOrderAgent>
            <Orders>
                <Order>
                	<TariffTypeId>1</TariffTypeId>
                	<Num>%s</Num>
                	<Date>%s</Date>
                	<Client>
                    	<FirstName>%s</FirstName>
                    	<LastName>%s</LastName>
                    	<MiddleName>%s</MiddleName>
                    	<ContactCellPhone>%s</ContactCellPhone>
                	</Client>
                	<Note>Техно выгоды. Интернет, ТВ, СВЯЗЬ Продавец - ИП Кривошеин ЯП</Note>
                	%s
                	<InstAdr>
                    	<RegionId>%s</RegionId>
                    	<CityId>%s</CityId>
                    	<StreetId>%s</StreetId>
                    	<HouseId>%s</HouseId>
                    	<Flat>%s</Flat>
                    	<TypeAdrId>0</TypeAdrId>
                	</InstAdr>
                </Order>
            </Orders>
        </CreatePacketOrderAgent>`, idLead, dateRequest, name, surname, patronymic, phone, xmlProduct, idRegion, idDistrict, idStreet, idHouse, flat)

	fmt.Println(requestBody)
	cert, err := tls.LoadX509KeyPair(e.pathCert, e.pathKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке сертификата: %w", err)
	}

	client := getClient(cert)

	req, err := http.NewRequest("POST", e.urlCRM, strings.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Content-Type", "text/xml")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка: получен неожиданный статус ответа %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении тела ответа: %w", err)
	}
	fmt.Println(string(body))

	var data []models.ResponseSendLead 
	err = xml.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("ошибка при разборе XML: %w", err)
	}

	return data, nil
}
func (e *EISSD) GetTariff(region string) (string, error) {
	tarrifs, err := e.repo.GetTariffsByRegion("01")
	if err != nil {
		return "", err
	}

	var tariffId int
	var techId int
	var techSvcId int
	var optionId int
	var optionSvcId int
	for _, tariff := range tarrifs {
		if tariff.Techs != nil && tariff.Options != nil {
			tariffId = tariff.ID
			techId = tariff.Techs[0].TechId
			techSvcId = tariff.Techs[0].SvcClassId
			optionId = tariff.Options[0].ID
			optionSvcId = tariff.Options[0].SvcClassId
		}
	}

	result := fmt.Sprintf(`
		<Product>
            <TariffId>%d</TariffId>
                <Techs>
                    <Tech>
                        <TechId>%d</TechId>
                        <serviceId>%d</serviceId>
                    </Tech>
                </Techs>
            <options>
                <option>
                    <OptionId>%d</OptionId>
                    <serviceId>%d</serviceId>
                </option>
            </options>
        </Product>`, tariffId, techId, techSvcId, optionId, optionSvcId)
	return result, nil
}