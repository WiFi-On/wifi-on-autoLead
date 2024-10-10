package models

type AddressDadata struct {
	RegionID string `json:"region_kladr_id"`

	RegionWithType string `json:"region_with_type"`
	RegionType string	`json:"region_type"`
    RegionTypeFull string  `json:"region_type_full"`
	Region string  `json:"region"`

	CityFiasId string `json:"city_fias_id"`
	CityWithType string `json:"city_with_type"`
	CityType string	`json:"city_type"`
    CityTypeFull string  `json:"city_type_full"`
	City string  `json:"city"`   

	SettlementFiasId string `json:"settlement_fias_id"`
	SettlementWithType string `json:"settlement_with_type"`
	SettlementType string	`json:"settlement_type"`
    SettlementTypeFull string  `json:"settlement_type_full"`
	Settlement string  `json:"settlement"` 

    StreetWithType string `json:"street_with_type"`
	StreetType string	`json:"street_type"`
    StreetTypeFull string  `json:"street_type_full"`
	Street string  `json:"street"`

	AreaWithType string `json:"area_with_type"`
	AreaType string	`json:"area_type"`
    AreaTypeFull string  `json:"area_type_full"`
	Area string  `json:"area"`

	HouseType string `json:"house_type"`
	HouseTypeFull string `json:"house_type_full"`
	House string `json:"house"`

	FlatType string `json:"flat_type"`
	FlatTypeFull string `json:"flat_type_full"`
	Flat string `json:"flat"`

	BlockType string `json:"block_type"`
	BlockTypeFull string `json:"block_type_full"`
	Block string `json:"block"`
}

type AddressResponseDadata struct {
	Suggestions []struct {
		Value             string   `json:"value"`
		UnrestrictedValue string   `json:"unrestricted_value"`
		Data              *AddressDadata `json:"data"`
	} `json:"suggestions"`
}



