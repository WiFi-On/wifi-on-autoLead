package models

type ResponseTariffInRepository struct {
    ID      int     `json:"id"`
    Name    string  `json:"name"`
    Region  string     `json:"region"`
    Techs   []struct {
        TechId     int `json:"TechId"`
        SvcClassId int `json:"SvcClassId"`
    } `json:"techs"`
    Cities  []struct {
        ID    int `json:"id"`
        Allow int `json:"allow"`
    } `json:"cities"`
    Options  []struct {
        ID                int    `json:"Id"`
        Fee               int    `json:"Fee"`
        Cost              int    `json:"Cost"`
        Name              string `json:"Name"`
        IsBase            int    `json:"IsBase"`
        TechId            int    `json:"TechId"`
        SvcClassId        int    `json:"SvcClassId"`
        SpeedKbPerSec     int    `json:"SpeedKbPerSec"`
        SalesChannelsId   string `json:"SalesChannelsId"`
    } `json:"params"`
}