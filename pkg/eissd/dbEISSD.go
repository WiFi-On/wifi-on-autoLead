package eissd

import (
	"database/sql"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

type DB struct {
	Conn *sql.DB
}

// NewDB открывает соединение с базой данных PostgreSQL
func NewDB(dataSourceName string) (*DB, error) {
	conn, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Проверяем соединение с базой данных
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &DB{Conn: conn}, nil
}

// Close закрывает соединение с базой данных
func (db *DB) Close() error {
	return db.Conn.Close()
}

// CreateDistrictsTable создает таблицу districts
func (db *DB) CreateDistrictsTable() error {
	query := `CREATE TABLE IF NOT EXISTS districts (
		id INTEGER PRIMARY KEY,
		region INTEGER NOT NULL,
		name TEXT NOT NULL,
		object TEXT NOT NULL
	);`
	_, err := db.Conn.Exec(query)
	return err
}
// CreateStreetsTable создает таблицу streets
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
// CreateHousesTable создает таблицу houses
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
// Добавляет информацию о населенном пункте в таблицу districts
func (db *DB) AddDistrict(id int, region int, name string, object string) error {
	query := "INSERT INTO districts (id, region, name, object) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING"

	_, err := db.Conn.Exec(query, id, region, name, object)
	if err != nil {
		return err
	}

	return nil
}
// Добавляет информацию о улице в таблицу streets
func (db *DB) AddStreet(id int, region string, name string, object string, districtID int) error {
	query := "INSERT INTO streets (id, region, name, object, district_id) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING"

	_, err := db.Conn.Exec(query, id, region, name, object, districtID)
	if err != nil {
		return err
	}

	return nil
}
// Добавляет информацию о доме в таблицу houses
func (db *DB) AddHouse(id string, region string, house string, streetID int) error {
	query := "INSERT INTO houses (id, region, house, street_id) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING"

	_, err := db.Conn.Exec(query, id, region, house, streetID)
	if err != nil {
		return err
	}

	return nil
}
