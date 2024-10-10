package models

type BitrixLead struct {
    ID string `json:"ID"`
    Title string `json:"TITLE"`
    StageID string `json:"STAGE_ID"`
    Opportunity string `json:"OPPORTUNITY"`
    ProviderID string `json:"UF_CRM_1697294773665"`
    TariffName string `json:"UF_CRM_1697294796468"`
    Comment string `json:"UF_CRM_1697294923031"`
    Number string `json:"UF_CRM_1697365970828"`
    Address string `json:"UF_CRM_1697646751446"`
    Client string `json:"UF_CRM_1697357613372"`
}

type BitrixLeadsResponse struct {
    Result []BitrixLead `json:"result"`
}