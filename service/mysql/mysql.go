package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"pdrygala.com/weather/core/weather"
)

type MySQL struct {
	conn *sql.DB
}

func NewMySQL(dsn string) (*MySQL, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %s: %w", dsn, err)
	}
	return &MySQL{
		conn: conn,
	}, nil
}

func (m *MySQL) Close(ctx context.Context) error {
	return m.conn.Close()
}

func ConnectDB(username string, password string, host string, port string, dbName string) (*sql.DB, error) {

	// Create a DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName)

	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", dsn)

	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTable(db *sql.DB) {
	createTableQuery := `
	CREATE TABLE weather (
	id INT AUTO_INCREMENT PRIMARY KEY,
	location VARCHAR(255) NOT NULL,
	date_time DATETIME NOT NULL,
	weather_code INT,     
	description VARCHAR(255),     
	temperature DECIMAL(5, 2),     
	wind_speed DECIMAL(5, 2),     
	wind_direction VARCHAR(10)
	)`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		panic(err.Error())
	}
}

func (m *MySQL) GetLatestRecord() string {

	// Perform the SQL SELECT query
	query := "SELECT * FROM weather ORDER BY id DESC LIMIT 1"
	var result weather.WeatherData

	err := m.conn.QueryRow(query).Scan(&result.Id, &result.City, &result.FormattedTime, &result.WeatherCode, &result.Description, &result.Temperature, &result.WindSpeed, &result.Direction)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No records found.")
		} else {
			panic(err.Error())
		}
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}
	// Print the JSON data as a string
	return (string(jsonData))
}

func GetAllRecords() string {
	// Perform connection to database
	db, err := ConnectDB("foobar", "password", "127.0.0.1", "3306", "db")

	//If there is an error handle it
	if err != nil {
		panic(err.Error())
	}
	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	query := "SELECT * FROM weather"
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var results []weather.WeatherData

	for rows.Next() {
		var result weather.WeatherData
		err := rows.Scan(&result.Id, &result.City, &result.FormattedTime, &result.WeatherCode, &result.Description, &result.Temperature, &result.WindSpeed, &result.Direction)
		if err != nil {
			panic(err.Error())
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		panic(err.Error())
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		panic(err.Error())
	}
	return (string(jsonData))
}

func GetResultsByTimeRange(startTime, endTime time.Time) string {
	// Perform connection to database
	db, err := ConnectDB("foobar", "password", "127.0.0.1", "3306", "db")

	//If there is an error handle it
	if err != nil {
		panic(err.Error())
	}
	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM weather WHERE date_time BETWEEN '%v' AND '%v'", startTime, endTime)

	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var results []weather.WeatherData

	for rows.Next() {
		var result weather.WeatherData
		err := rows.Scan(&result.Id, &result.City, &result.FormattedTime, &result.WeatherCode, &result.Description, &result.Temperature, &result.WindSpeed, &result.Direction)
		if err != nil {
			panic(err.Error())
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		panic(err.Error())
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		panic(err.Error())
	}
	return (string(jsonData))
}
