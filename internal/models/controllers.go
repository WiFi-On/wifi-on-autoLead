package models

// AddressData содержит адрес одной строчкой. Используется для получения и для отправки.
type AddressData struct {
	Address string `json:"address"`
}

// ErrorResponse содержит информацию об ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}

type CheckTHVResponse struct {
	Success bool   `json:"success"`
}
