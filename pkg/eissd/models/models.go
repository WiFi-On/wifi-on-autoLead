package eissd

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
	RegionId string `xml:"RegionId"`
	StreetId int `xml:"StreetId"`
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
