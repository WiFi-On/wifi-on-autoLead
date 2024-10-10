package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"wifionAutolead/internal/models"
)

type Bitrix struct {
	URLToken string
}

func NewBitrix(urltoken string) *Bitrix {
	return &Bitrix{
		URLToken: urltoken,
	}
}

// id ростелекома из битрикса = 52
// Функция получения списка лидов по айди(из битрикса) провайдера.
func (b *Bitrix) GetDealsOnProviders(idProvider string) ([]models.BitrixLead, error) {
	if b.URLToken == "" {
		return nil, fmt.Errorf("URLToken is empty")
	}

	var urlGetLeads string = fmt.Sprintf("%s%s", b.URLToken, "crm.deal.list?")

	params := url.Values{}
	params.Add("filter[UF_CRM_1697294773665]", idProvider)
	params.Add("filter[STAGE_ID]", "PREPAYMENT_INVOICE")
	params.Add("select[]", "*")
	params.Add("select[]", "UF_*")

	urlGetLeads = fmt.Sprintf("%s%s", urlGetLeads, params.Encode())
	fmt.Printf("URL: %s\n", urlGetLeads)

	resp, err := http.Get(urlGetLeads)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	var response models.BitrixLeadsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании JSON: %v", err)
	}

	return response.Result, nil
}

