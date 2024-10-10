package eissd

import "encoding/xml"

// Address содержит информацию об одном адресе из XML-ответа
type Address struct {
	RegionId       string `xml:"RegionId"`
	AddrObjectId   int    `xml:"AddrObjectId"`
	NameAddrObject string `xml:"NameAddrObject"`
	AbbrNameObject string `xml:"AbbrNameObject"`
	ParentId       int    `xml:"ParentId"`
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
	XMLName  xml.Name `xml:"GetTariffPlansAgent"`
	Tariffs  Tariffs  `xml:"tariffs"`
	Response string   `xml:"response"`
	Message  string   `xml:"message"`
}
// Tariffs представляет список тарифов
type Tariffs struct {
	Tariff []Tariff `xml:"Tariff"`
}
// Tariff представляет тариф
type Tariff struct {
	PublicName        string  `xml:"PublicName"`
	TariffId          int     `xml:"TariffId"`
	SalesChannelsId   string  `xml:"SalesChannelsId"`
	Techs             Techs   `xml:"Techs"`
	TVChannelCount    int     `xml:"TVChannelCount"`
	IsPromo           int     `xml:"IsPromo"`
	RegionId          int     `xml:"RegionId"`
	TrafficType       int     `xml:"TrafficType"`
	MarketInfo        string  `xml:"MarketInfo"`
	DeviceRequirement int     `xml:"DeviceRequirement"`
	Cities            Cities  `xml:"Cities"`
	SvcClassIds       []int   `xml:"SvcClassIds>SvcClassId"`
	Options           Options `xml:"options"`
	TypeId            int     `xml:"typeId"`
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
	SvcClassId      int    `xml:"SvcClassId"`
	ID              int    `xml:"Id"`
	Name            string `xml:"Name"`
	TechId          int    `xml:"TechId"`
	Cost            int    `xml:"Cost"`
	Fee             int    `xml:"Fee"`
	SpeedKbPerSec   int    `xml:"SpeedKbPerSec"`
	IsBase          int    `xml:"IsBase"`
	SalesChannelsId string `xml:"SalesChannelsId"`
}
type GetTariffPlansMvno struct {
	XMLName  xml.Name `xml:"GetTariffPlansMvno"`
	Response int      `xml:"Response"`
	Message  string   `xml:"Message"`
	Tariffs  struct {
		Tariff []struct {
			TarId        int    `xml:"TarId"`
			PSTarId      int    `xml:"PSTarId"`
			RegionMvnoId string    `xml:"RegionMvnoId"`
			StartDate    string `xml:"StartDate"`
			Title        string `xml:"Title"`
			Category     int    `xml:"Category"`
			IsForDealer  string `xml:"IsForDealer"`
			IsActive     int    `xml:"IsActive"`
			TarType      int    `xml:"TarType"`
			PayType      int    `xml:"PayType"`
			NumberTypes  struct {
				NumberType struct {
					PhoneFederal int     `xml:"PhoneFederal"`
					PhoneColorId int     `xml:"PhoneColorId"`
					Cost         float64 `xml:"Cost"`
					Advance      float64 `xml:"Advance"`
				} `xml:"NumberType"`
			} `xml:"NumberTypes"`
		} `xml:"Tariff"`
	} `xml:"Tariffs"`
}
type GetTariffPlansMvnoJSON struct {
	XMLName  xml.Name `json:"GetTariffPlansMvno"`
	Response int      `json:"Response"`
	Message  string   `json:"Message"`
	Tariffs  struct {
		Tariff []struct {
			TarId        int    `json:"TarId"`
			PSTarId      int    `json:"PSTarId"`
			RegionMvnoId string    `json:"RegionMvnoId"`
			StartDate    string `json:"StartDate"`
			Title        string `json:"Title"`
			Category     int    `json:"Category"`
			IsForDealer  string `json:"IsForDealer"`
			IsActive     int    `json:"IsActive"`
			TarType      int    `json:"TarType"`
			PayType      int    `json:"PayType"`
			NumberTypes  struct {
				NumberType struct {
					PhoneFederal int     `json:"PhoneFederal"`
					PhoneColorId int     `json:"PhoneColorId"`
					Cost         float64 `json:"Cost"`
					Advance      float64 `json:"Advance"`
				} `json:"NumberType"`
			} `json:"NumberTypes"`
		} `json:"Tariff"`
	} `json:"Tariffs"`
}
// Модели для работы с тарифами в формате JSON
// type GetTariffPlansAgent struct {
// 	Tariffs  Tariffs `json:"tariffs"`
// 	Response string  `json:"response"`
// 	Message  string  `json:"message"`
// }
// type Tariffs struct {
// 	Tariff []Tariff `json:"Tariff"`
// }
// type Tariff struct {
// 	PublicName        string  `json:"PublicName"`
// 	TariffId          int     `json:"TariffId"`
// 	SalesChannelsId   string  `json:"SalesChannelsId"`
// 	Techs             Techs   `json:"Techs"`
// 	TVChannelCount    int     `json:"TVChannelCount"`
// 	IsPromo           int     `json:"IsPromo"`
// 	RegionId          int     `json:"RegionId"`
// 	TrafficType       int     `json:"TrafficType"`
// 	MarketInfo        string  `json:"MarketInfo"`
// 	DeviceRequirement int     `json:"DeviceRequirement"`
// 	Cities            Cities  `json:"Cities"`
// 	SvcClassIds       []int   `json:"SvcClassIds"`
// 	Options           Options `json:"options"`
// 	TypeId            int     `json:"typeId"`
// }
// type Techs struct {
// 	Tech []Tech `json:"Tech"`
// }
// type Tech struct {
// 	TechId     int `json:"TechId"`
// 	SvcClassId int `json:"SvcClassId"`
// }
// type Cities struct {
// 	City []City `json:"city"`
// }
// type City struct {
// 	ID    int `json:"id"`
// 	Allow int `json:"allow"`
// }
// type Options struct {
// 	Option []Option `json:"option"`
// }
// type Option struct {
// 	SvcClassId      int    `json:"SvcClassId"`
// 	ID              int    `json:"Id"`
// 	Name            string `json:"Name"`
// 	TechId          int    `json:"TechId"`
// 	Cost            int    `json:"Cost"`
// 	Fee             int    `json:"Fee"`
// 	SpeedKbPerSec   int    `json:"SpeedKbPerSec"`
// 	IsBase          int    `json:"IsBase"`
// 	SalesChannelsId string `json:"SalesChannelsId"`
// }