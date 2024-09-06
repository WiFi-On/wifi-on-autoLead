package eissd

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Conn *sql.DB
}

// NewDB открывает соединение с базой данных
func NewDB(dataSourceName string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DB{Conn: conn}, nil
}

// Close закрывает соединение с базой данных
func (db *DB) Close() error {
	return db.Conn.Close()
}

// CreateUserTable создает таблицу users
func (db *DB) CreateDistrictsTable() error {
	query := `CREATE TABLE IF NOT EXISTS districts (
		id INTEGER PRIMARY KEY not null,
		region INTEGER not null,
		name TEXT not null,
		object TEXT not null
	);`
	_, err := db.Conn.Exec(query)
	return err
}

func (db *DB) CreateStreetsTable() error {
	query := `CREATE TABLE IF NOT EXISTS streets (
		id INTEGER PRIMARY KEY,
		region TEXT,
		name TEXT,
		object TEXT,
		district_id INTEGER REFERENCES districts(id)
	);`
	_, err := db.Conn.Exec(query)
	return err
}

func (db *DB) CreateHousesTable() error {
	query := `CREATE TABLE IF NOT EXISTS houses (
		id INTEGER PRIMARY KEY,
		region TEXT,
		house TEXT,
		street_id INTEGER REFERENCES streets(id)
	);`
	_, err := db.Conn.Exec(query)
	return err
}

func (db *DB) AddDistrict(id int, region int, name string, object string) error {
	query := "INSERT INTO districts (id, region, name, object) VALUES (?, ?, ?, ?)"

	_, err := db.Conn.Exec(query, id, region, name, object)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) AddStreet(id int, region string, name string, object string, districtID int) error {
	query := "INSERT INTO streets (id, region, name, object, district_id) VALUES (?, ?, ?, ?, ?)"

	_, err := db.Conn.Exec(query, id, region, name, object, districtID)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) AddHouses(id int, region string, house string, streetID int) error {
	query := "INSERT INTO houses (id, region, house, street_id) VALUES (?, ?, ?, ?)"

	_, err := db.Conn.Exec(query, id, region, house, streetID)
	if err != nil {
		return err
	}

	return nil
}

// Получение id района (district) по region и name
func (db *DB) GetDistrictIDByRegionAndName(region int, name string) (int, error) {
	var districtID int

	query := "SELECT id FROM districts WHERE region = ? AND name = ?"
	err := db.Conn.QueryRow(query, region, name).Scan(&districtID)
	if err != nil {
		return 0, err
	}

	return districtID, nil
}

// Получение id улицы (street) по region, name и districtID
func (db *DB) GetStreetIDByRegionNameAndDistrict(region int, name string, districtID int) (int, error) {
	var streetID int

	query := "SELECT id FROM streets WHERE region = ? AND name = ? AND district_id = ?"
	err := db.Conn.QueryRow(query, region, name, districtID).Scan(&streetID)
	if err != nil {
		return 0, err
	}

	return streetID, nil
}

// Получение id дома по region, streetID и house
func (db *DB) GetHouseIDByRegionStreetAndHouse(region int, streetID int, house int) (int, error) {
	var houseID int

	query := "SELECT id FROM houses WHERE region = ? AND street_id = ? AND house = ?"
	err := db.Conn.QueryRow(query, region, streetID, house).Scan(&houseID)
	if err != nil {
		return 0, err
	}

	return houseID, nil
}
