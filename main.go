package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Response struct {
	CurrentWeather struct {
		Temperature   float64 `json:"temperature"`
		WindSpeed     float64 `json:"windspeed"`
		WindDirection float64 `json:"winddirection"`
		WeatherCode   int     `json:"weathercode"`
		IsDay         int     `json:"is_day"`
		Time          string  `json:"time"`
	} `json:"current_weather"`
}

var weatherCodeMap = map[int]string{
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

type WindDirection struct {
	DegreeStart int
	DegreeEnd   int
	Description string
}

var windDirections = []WindDirection{
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

func getWindDirection(degree int) string {
	for _, direction := range windDirections {
		if degree >= direction.DegreeStart && degree < direction.DegreeEnd {
			return direction.Description
		}
	}
	return "Unknown"
}

func parseDateTime(dateTimeStr string) (time.Time, error) {
	// Parse the string into a time.Time value
	layout := "2006-01-02T15:04"
	dateTime, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		return time.Time{}, err
	}

	return dateTime, nil
}

func convertToWarsawTime(dateTime time.Time) (time.Time, error) {
	// Convert to Warsaw time (UTC+1)
	location, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		return time.Time{}, err
	}
	warsawTime := dateTime.In(location)

	return warsawTime, nil
}

func formatDateTime(dateTime time.Time) string {
	// Format the time in a easy-readable format
	formattedTime := dateTime.Format("2006-01-02 15:04 MST")
	return formattedTime
}

func main() {
	url := "https://api.open-meteo.com/v1/forecast?latitude=54.52&longitude=18.53&current_weather=true"

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error making Get request:", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var jsonBody Response

	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	dateTime, err := parseDateTime(jsonBody.CurrentWeather.Time)
	if err != nil {
		fmt.Println("Error parsing date and time:", err)
		return
	}

	warsawTime, err := convertToWarsawTime(dateTime)
	if err != nil {
		fmt.Println("Error converting time to Warsow", err)
		return
	}
	formattedTime := formatDateTime(warsawTime)

	temperature := jsonBody.CurrentWeather.Temperature
	windSpeed := jsonBody.CurrentWeather.WindSpeed
	weatherCode := jsonBody.CurrentWeather.WeatherCode
	direction := getWindDirection(int(jsonBody.CurrentWeather.WindDirection))

	fmt.Println("Lokalizacja: Gdynia")
	fmt.Printf("Data pomiaru: %s \n", formattedTime)
	description, exists := weatherCodeMap[weatherCode]
	if exists {
		fmt.Printf("Aktualna pogoda: %s\n", description)
	} else {
		fmt.Printf("Nieznany opis dla kodu pogody: %d\n", weatherCode)
	}
	fmt.Println("Temperatura:", temperature)
	fmt.Println("Wiatr:", windSpeed)
	fmt.Printf("Kierunek wiatru: %s\n", direction)

}
