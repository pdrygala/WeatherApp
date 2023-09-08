package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"pdrygala.com/weather/core/weather"
)

func GetLatestRecord() string {

	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", "foobar:password@tcp(127.0.0.1:3306)/db")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	// Perform the SQL SELECT query
	query := "SELECT * FROM weather_data ORDER BY id DESC LIMIT 1"
	var result weather.WeatherDB // Replace with your actual struct type

	// ///	FormattedTime string
	// City          string
	// WeatherCode   int
	// Description   string
	// Temperature   string
	// WindSpeed     string
	// Direction     string

	err = db.QueryRow(query).Scan(&result.Id, &result.City, &result.FormattedTime, &result.Description, &result.Temperature, &result.WindSpeed, &result.Direction)
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
