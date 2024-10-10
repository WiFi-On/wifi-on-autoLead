package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"wifionAutolead/internal/models"

	_ "github.com/lib/pq" // Импортируйте драйвер PostgreSQL
)

type DB struct {
    Conn *sql.DB
}

// Создание нового экземпляра DB с подключением к базе данных
func NewDB(dataSourceName string) (*DB, error) {
    db, err := sql.Open("postgres", dataSourceName) // Замените "postgres" на нужный драйвер
    if err != nil {
        return nil, err
    }

    // Дополнительно проверьте подключение к базе данных
    if err := db.Ping(); err != nil {
        return nil, err
    }

    return &DB{Conn: db}, nil
}
// Получение id района (district) по region и name
func (db *DB) GetDistrictIDByRegionAndName(region string, name string) (string,  error) {
    var districtID string

    query := "SELECT id FROM districts WHERE region = $1 AND name = $2"
    err := db.Conn.QueryRow(query, region, name).Scan(&districtID)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", nil // если данные не найдены
        }
        return "",  err
    }

    return districtID,  nil
}
// Получение id района (district) по parentID
func (db *DB) GetDistrictIDByParentIDandName(parentID string, name string) (string, error) {
    var districtID string

    query := "SELECT id FROM districts WHERE parent_id = $1 and name = $2"
    err := db.Conn.QueryRow(query, parentID, name).Scan(&districtID)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", nil // если данные не найдены
        }
        return "",  err
    }

    return districtID,  nil
}
// Получение id родительского района (district) по region и name
func (db *DB) GetParentIDByRegionAndName(region string, name string) (string, error) {
    var parentID string

    query := "SELECT parent_id FROM districts WHERE region = $1 AND name = $2"
    err := db.Conn.QueryRow(query, region, name).Scan(&parentID)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", nil // если данные не найдены
        }
        return "", err
    }

    return parentID, nil
}
// Получение id родительского района (district) по region, name и parentID
func (db *DB) GetParentIDByDistrictID(districtID string) (string, error) {
    var parentID string

    query := "SELECT parent_id FROM districts WHERE id = $1 "
    err := db.Conn.QueryRow(query, districtID).Scan(&parentID)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", nil // если данные не найдены
        }
        return "", err
    }

    return parentID, nil
}
// Получение id улицы (street) по region, name и districtID
func (db *DB) GetStreetIDByNameAndDistrictId(name string, districtID string) (string, error) {
    var streetID string

    query := "SELECT id FROM streets WHERE name = $1 AND district_id = $2"
    err := db.Conn.QueryRow(query, name, districtID).Scan(&streetID)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", nil // если данные не найдены
        }
        return "", err
    }

    return streetID, nil
}
// Получение id дома по region, streetID и house
func (db *DB) GetHouseIDByStreetIdAndHouse( streetID string, house string) (string, error) {
    var houseID string

    query := "SELECT id FROM houses WHERE street_id = $1 AND house = $2"
    err := db.Conn.QueryRow(query, streetID, house).Scan(&houseID)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", nil // если данные не найдены
        }
        return "", err
    }

    return houseID, nil
}
// Получение тарифов по региону
func (db *DB) GetTariffsByRegion(region string) ([]models.ResponseTariffInRepository, error) {
    rows, err := db.Conn.Query("SELECT id, name, region, techs, cities, options FROM tariffs WHERE region = $1", region)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tariffs []models.ResponseTariffInRepository
    for rows.Next() {
        var (
            id      int
            name    string
            region  string
            techs   string
            cities  string
            options  string
        )

        if err := rows.Scan(&id, &name, &region, &techs, &cities, &options); err != nil {
            return nil, err
        }

        var techsSlice []struct {
            TechId     int `json:"TechId"`
            SvcClassId int `json:"SvcClassId"`
        }
        var citiesSlice []struct {
            ID    int `json:"id"`
            Allow int `json:"allow"`
        }
        var optionsSlice []struct {
            ID                int    `json:"Id"`
            Fee               int    `json:"Fee"`
            Cost              int    `json:"Cost"`
            Name              string `json:"Name"`
            IsBase            int    `json:"IsBase"`
            TechId            int    `json:"TechId"`
            SvcClassId        int    `json:"SvcClassId"`
            SpeedKbPerSec     int    `json:"SpeedKbPerSec"`
            SalesChannelsId   string `json:"SalesChannelsId"`
        }

        if err := json.Unmarshal([]byte(techs), &techsSlice); err != nil {
            return nil, err
        }
        if err := json.Unmarshal([]byte(cities), &citiesSlice); err != nil {
            return nil, err
        }
        if err := json.Unmarshal([]byte(options), &optionsSlice); err != nil {
            return nil, err
        }

        tariff := models.ResponseTariffInRepository{
            ID:     id,
            Name:   name,
            Region: region,
            Techs:  techsSlice,
            Cities: citiesSlice,
            Options: optionsSlice,
        }

        tariffs = append(tariffs, tariff)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    fmt.Println(tariffs)

    return tariffs, nil
}
