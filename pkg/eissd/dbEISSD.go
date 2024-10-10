package eissd

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DB struct {
	Conn *sql.DB
}
// NewDB открывает соединение с базой данных
func NewDB(dataSourceName string) (*DB, error) {
	conn, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DB{Conn: conn}, nil
}
// Close закрывает соединение с базой данных
func (db *DB) Close() error {
	return db.Conn.Close()
}
// CreateUserTable создает таблицу Districts в базе данных
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
// CreateStreetsTable создает таблицу Streets в базе данных
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
// CreateHousesTable создает таблицу Houses в базе данных
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
// CreateTariffsTable создает таблицу Tariffs в базе данных
func (db *DB) CreateTariffsTable() error {
	query := `CREATE TABLE IF NOT EXISTS tariffs (
		id INTEGER PRIMARY KEY,
		name TEXT,
		region TEXT,
		techs jsonb,
		cities jsonb,
		options jsonb
	)`
	_, err := db.Conn.Exec(query)

	return err
}
// CreateTariffsMVNOTable создает таблицу TariffsMVNO в базе данных
func (db *DB) CreateTariffsMVNOTable() error {
	query := `create table tariffsMVNO(
    tar_id int not null,
    pst_tar_id INTEGER NOT NULL,
    region_mvno_id text NOT NULL,
    start_date text NOT NULL,
    title VARCHAR(255) NOT NULL,
    category INTEGER NOT NULL,
    is_for_dealer text NOT NULL,
    is_active int NOT NULL,
    tar_type int NOT NULL,
    pay_type int NOT NULL,
    number_types jsonb not NULL
	);`
	_, err := db.Conn.Exec(query)

	return err
}
// CreateTariffsClientTable создает таблицу TariffsClient в базе данных
func (db *DB) CreateTariffsClientTable() error {
	query := `CREATE TABLE IF NOT EXISTS tariffsClient (
		id INTEGER PRIMARY KEY
		district_id INTEGER REFERENCES districts(id)
		internet_speed INTEGER
		channels_count INTEGER
		minutes INTEGER
		gigabytes INTEGER
		sms INTEGER
		connection_cost INTEGER
		cost INTEGER
		sale_description TEXT
		name TEXT
		router_rent INTEGER
		router_cost INTEGER
		router_payment INTEGER
		tv_box_payment INTEGER
		technologies jsonb
		type INTEGER
		additional_info jsonb
	)`
	_, err := db.Conn.Exec(query)

	return err
}
// AddDistrict добавляет новый дистрикт в базу данных
func (db *DB) AddDistrict(id int, region string, name string, object string, parentID int) error {
	query := "INSERT INTO districts (id, region, name, object, parent_id) VALUES ($1, $2, $3, $4, $5)"

	_, err := db.Conn.Exec(query, id, region, name, object, parentID)
	if err != nil {
		return err
	}

	return nil
}
// AddStreet добавляет новую улицу в базу данных
func (db *DB) AddStreet(id int, region string, name string, object string, districtID int) error {
	query := "INSERT INTO streets (id, region, name, object, district_id) VALUES ($1, $2, $3, $4, $5)"

	_, err := db.Conn.Exec(query, id, region, name, object, districtID)
	if err != nil {
		return err
	}

	return nil
}
// AddHouse добавляет новый дом в базу данных
func (db *DB) AddHouses(id int, region string, house string, streetID int) error {
	query := "INSERT INTO houses (id, region, house, street_id) VALUES ($1, $2, $3, $4)"

	_, err := db.Conn.Exec(query, id, region, house, streetID)
	if err != nil {
		return err
	}

	return nil
}
// AddTariff добавляет новый тариф в базу данных
func (db *DB) AddTariff(id int, name string, region string, techsJSON string, citiesJSON string, optionsJSON string) error {
    query := `INSERT INTO tariffs (id, name, region, techs, cities, options) VALUES ($1, $2, $3, $4, $5, $6)`

    _, err := db.Conn.Exec(query, id, name, region, techsJSON, citiesJSON, optionsJSON)
    if err != nil {
        return err
    }

    return nil
}
// AddTariffMVNO добавляет новый тариф в базу данных
func (db *DB) AddTariffMVNO(tariffId int, PStatiffId int, regionMVNO string, startDate string, name string, category int, isForDealer string, isActive int, tarType int, patType int, numberTypesJSON string) error {
	query := `INSERT INTO tariffsMVNO (tar_id, pst_tar_id, region_mvno_id, start_date, title, category, is_for_dealer, is_active, tar_type, pay_type, number_types) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := db.Conn.Exec(query, tariffId, PStatiffId, regionMVNO, startDate, name, category, isForDealer, isActive, tarType, patType, numberTypesJSON)
	if err != nil {
		return err
	}

	return nil
}
// Добавляет id тарифов для региона
func (db *DB) AddTariffForRegion(regionID string, shpdID int, tvID int, mvnoID int , mltID int) error {
	query := `INSERT INTO tariffs_for_region (region_id, shpd_tarriff_id, tv_tarriff_id, mvno_tarriff_id, mlt_tarriff_id) VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Conn.Exec(query, regionID, shpdID, tvID, mvnoID, mltID)
	if err != nil {
		return err
	}

	return nil
}
// Добавление тарифа в таблицу tariffsClient

