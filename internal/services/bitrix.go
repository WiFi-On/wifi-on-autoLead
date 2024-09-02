package services

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Bitrix struct {
	URLToken string
}

func NewBitrix(urltoken string) *Bitrix {
	return &Bitrix{
		URLToken: urltoken,
	}
}

// Функция получения списка лидов по айди(из битрикса) провайдера.
func (b *Bitrix) GetDealsOnProviders(idProvider string) ([]string, error) {
	if b.URLToken == "" {
		return nil, fmt.Errorf("URLToken is empty")
	}

	var urlGetLeads string = fmt.Sprintf("%s%s", b.URLToken, "crm.deal.list?")

	params := url.Values{}
	params.Add("filter[UF_CRM_1697294773665]", idProvider)
	params.Add("filter[STAGE_ID]", "PREPAYMENT_INVOICE")
	


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

	return []string{string(body)}, nil
}
