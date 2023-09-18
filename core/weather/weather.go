package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// WeatherResponse represents the structure of the weather API response
type WeatherResponse struct {
	CurrentWeather struct {
		Temperature   float64 `json:"temperature"`
		WindSpeed     float64 `json:"windspeed"`
		WindDirection float64 `json:"winddirection"`
		WeatherCode   int     `json:"weathercode"`
		IsDay         int     `json:"is_day"`
		Time          string  `json:"time"`
	} `json:"current_weather"`
}

type WindDirection struct {
	DegreeStart int
	DegreeEnd   int
	Description string
}

var (
	weatherCodeMap = map[int]string{
		0:  "Bezchmurne niebo",
		1:  "Przeważnie bezchmurne",
		2:  "Częściowo zachmurzone",
		3:  "Zachmurzone",
		45: "Mgła i szadź",
		48: "Mgła i szadź",
		51: "Mżawka: Słaba intensywność",
		53: "Mżawka: Umiarkowana intensywność",
		55: "Mżawka: Gęsta intensywność",
		56: "Marznąca mżawka: Słaba intensywność",
		57: "Marznąca mżawka: Gęsta intensywność",
		61: "Deszcz: Słaba intensywność",
		63: "Deszcz: Umiarkowana intensywność",
		65: "Deszcz: Duża intensywność",
		66: "Marznący deszcz: Słaba intensywność",
		67: "Marznący deszcz: Duża intensywność",
		71: "Opady śniegu: Słaba intensywność",
		73: "Opady śniegu: Umiarkowana intensywność",
		75: "Opady śniegu: Duża intensywność",
		77: "Kryształki śniegu",
		80: "Przelotne opady deszczu: Słaba intensywność",
		81: "Przelotne opady deszczu: Umiarkowana intensywność",
		82: "Przelotne opady deszczu: Wysoka intensywność",
		85: "Przelotne opady śniegu: Słaba intensywność",
		86: "Przelotne opady śniegu: Duża intensywność",
		95: "Burza: Słaba intensywność",
		96: "Burza z drobnym gradem",
		99: "Burza z dużym gradem",
	}

	windDirections = []WindDirection{
		{DegreeStart: 335, DegreeEnd: 360, Description: "N"},
		{DegreeStart: 295, DegreeEnd: 335, Description: "NW"},
		{DegreeStart: 245, DegreeEnd: 295, Description: "W"},
		{DegreeStart: 205, DegreeEnd: 245, Description: "SW"},
		{DegreeStart: 155, DegreeEnd: 205, Description: "S"},
		{DegreeStart: 115, DegreeEnd: 155, Description: "SE"},
		{DegreeStart: 65, DegreeEnd: 115, Description: "E"},
		{DegreeStart: 25, DegreeEnd: 65, Description: "NE"},
		{DegreeStart: 0, DegreeEnd: 25, Description: "N"},
	}
)

func getWindDirection(degree int) string {
	for _, direction := range windDirections {
		if degree >= direction.DegreeStart && degree < direction.DegreeEnd {
			return direction.Description
		}
	}
	return "Unknown"
}

func HandleWeatherData(weatherData WeatherResponse, city string) {

	data := WeatherData{

		FormattedTime: formatDateTime(weatherData.CurrentWeather.Time),
		City:          city,
		Temperature:   fmt.Sprintf("%v", weatherData.CurrentWeather.Temperature),
		WindSpeed:     fmt.Sprintf("%v", weatherData.CurrentWeather.WindSpeed),
		WeatherCode:   weatherData.CurrentWeather.WeatherCode,
		Direction:     getWindDirection(int(weatherData.CurrentWeather.WindDirection)),
	}

	fmt.Printf("Lokalizacja: %s \n", city)
	fmt.Printf("Data pomiaru: %s \n", data.FormattedTime)
	description, exists := weatherCodeMap[weatherData.CurrentWeather.WeatherCode]
	if exists {
		data.Description = description
		fmt.Printf("Aktualna pogoda: %s\n", data.Description)
	} else {
		fmt.Printf("Nieznany opis dla kodu pogody: %d\n", weatherData.CurrentWeather.WeatherCode)
	}
	fmt.Println("Temperatura:", data.Temperature, "°C")
	fmt.Println("Wiatr:", data.WindSpeed, "km/h")
	fmt.Printf("Kierunek wiatru: %s\n", data.Direction)

	WriteToDatabase(data)

}

func FetchWeather(latitude, longitude string) (WeatherResponse, error) {

	host := "https://api.open-meteo.com/v1/forecast?"
	params := fmt.Sprintf("latitude=%s&longitude=%s&current_weather=true", latitude, longitude)
	url := host + params
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making Get request:", err)
	}
	defer resp.Body.Close()

	weatherBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading WeatherResponse body:", err)
		return WeatherResponse{}, err
	}
	var jsonWeatherBody WeatherResponse

	if err := json.Unmarshal(weatherBody, &jsonWeatherBody); err != nil {
		fmt.Println("Error parsing JSON for weather:", err)
		return WeatherResponse{}, err
	}
	return jsonWeatherBody, nil
}

func formatDateTime(dateTimeStr string) string {
	dateTime, err := time.Parse("2006-01-02T15:04", dateTimeStr)
	if err != nil {
		fmt.Println("Error parsing date and time:", err)
		return ""
	}

	// Convert to Warsaw time (UTC+1)
	location, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return ""
	}
	warsawTime := dateTime.In(location)

	// Format the time in an easy-readable format
	return warsawTime.Format("2006-01-02 15:04:00")
}

func WriteToDatabase(data WeatherData) {
	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", "foobar:password@tcp(127.0.0.1:3306)/db")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	// perform a db.Query insert
	// INSERT INTO weather_data VALUES (1,'Gdynia','2023-09-04 21:00:00','Przeważnie bezchmurne','18.2','4.5','SW');

	insert, err := db.Query("INSERT INTO weather (location, date_time, weather_code, description, temperature, wind_speed, wind_direction) VALUES (?, ?, ?, ?, ?, ?, ?)",
		data.City, data.FormattedTime, data.WeatherCode, data.Description, data.Temperature, data.WindSpeed, data.Direction)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()

}
