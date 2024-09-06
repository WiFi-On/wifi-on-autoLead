package eissd

import "encoding/xml"

// Address содержит информацию об одном адресе из XML-ответа
type Address struct {
	RegionId       string    `xml:"RegionId"`
	AddrObjectId   int `xml:"AddrObjectId"`
	NameAddrObject string `xml:"NameAddrObject"`
	AbbrNameObject string `xml:"AbbrNameObject"`
	ParentId       int `xml:"ParentId"`
}
// AddressHouse содержит информацию об одном доме из XML-ответа
type AddressHouse struct {
	RegionId int    `xml:"RegionId"`
	StreetId string `xml:"StreetId"`
	HouseId  string `xml:"HouseId"`
	House    string `xml:"House"`
}
// GetAddressInfoAgentResponse структура для обработки всего XML-ответа
type GetAddressInfoAgentResponse struct {
	Addresses []Address `xml:"addresss>address"`
}
// GetAddressHouseInfoAgentResponse структура для обработки всего XML-ответа
type GetAddressHouseInfoAgentResponse struct {
	AddressHouses []AddressHouse `xml:"AddressHouses>AddressHouse"`
}
// CheckConnectionPossibilityResponse описывает ответ на запрос проверки возможности подключения
type CheckConnectionPossibilityResponse struct {
	Response int    `xml:"Response"`
	Message  string `xml:"Message"`
}

// GetTariffPlansAgent представляет структуру для ответа XML
type GetTariffPlansAgent struct {
    XMLName   xml.Name `xml:"GetTariffPlansAgent"`
    Tariffs   Tariffs  `xml:"tariffs"`
    Response  string   `xml:"response"`
    Message   string   `xml:"message"`
}

// Tariffs представляет список тарифов
type Tariffs struct {
    Tariff []Tariff `xml:"Tariff"`
}

// Tariff представляет тариф
type Tariff struct {
    PublicName        string   `xml:"PublicName"`
    TariffId          int      `xml:"TariffId"`
    SalesChannelsId   string   `xml:"SalesChannelsId"`
    Techs             Techs    `xml:"Techs"`
    TVChannelCount    int      `xml:"TVChannelCount"`
    IsPromo           int      `xml:"IsPromo"`
    RegionId          int      `xml:"RegionId"`
    TrafficType       int      `xml:"TrafficType"`
    MarketInfo        string   `xml:"MarketInfo"`
    DeviceRequirement int     `xml:"DeviceRequirement"`
    Cities            Cities   `xml:"Cities"`
    SvcClassIds       []int    `xml:"SvcClassIds>SvcClassId"`
    Options           Options  `xml:"options"`
    TypeId            int      `xml:"typeId"`
}

// Techs представляет список технологий
type Techs struct {
    Tech []Tech `xml:"Tech"`
}

// Tech представляет технологию
type Tech struct {
    TechId     int `xml:"TechId"`
    SvcClassId int `xml:"SvcClassId"`
}

// Cities представляет список городов
type Cities struct {
    City []City `xml:"city"`
}

// City представляет город
type City struct {
    ID    int `xml:"id"`
    Allow int `xml:"allow"`
}

// Options представляет список опций
type Options struct {
    Option []Option `xml:"option"`
}

// Option представляет опцию
type Option struct {
    SvcClassId       int    `xml:"SvcClassId"`
    ID               int    `xml:"Id"`
    Name             string `xml:"Name"`
    TechId           int    `xml:"TechId"`
    Cost             int    `xml:"Cost"`
    Fee              int    `xml:"Fee"`
    SpeedKbPerSec    int    `xml:"SpeedKbPerSec"`
    IsBase           int    `xml:"IsBase"`
    SalesChannelsId  string `xml:"SalesChannelsId"`
}

