package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

// CityResponse represents the structure of the city API response
type CityResponse struct {
	Address struct {
		HouseNumber  string `json:"house_number"`
		Road         string `json:"road"`
		Suburb       string `json:"suburb"`
		CityDistrict string `json:"city_district"`
		City         string `json:"city"`
		State        string `json:"state"`
		Postcode     string `json:"postcode"`
		Country      string `json:"country"`
		CountryCode  string `json:"country_code"`
	} `json:"address"`
}

// WindDirection represents a wind direction range and description

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

func getCity(latitude string, longitude string) string {
	// Get the city from the coordinates
	url := fmt.Sprintf("https://geocode.maps.co/reverse?lat=%s&lon=%s", latitude, longitude)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making Get request:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading WeatherResponse body:", err)
	}

	var jsonCityBody CityResponse

	if err := json.Unmarshal(body, &jsonCityBody); err != nil {
		fmt.Println("Error parsing JSON for city:", err)
		return err.Error()
	}

	return jsonCityBody.Address.City
}

func writeToFile(fileName string, data []string) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(data); err != nil {
		panic(err)
	}
}

func fetchWeather(latitude, longitude string) (WeatherResponse, error) {

	host := "https://api.open-meteo.com/v1/forecast?"
	params := fmt.Sprintf("latitude=%s&longitude=%s&current_weather=true", latitude, longitude)
	url := host + params

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

func handleWeatherData(weatherData WeatherResponse, city string) {

	formattedTime := formatDateTime(weatherData.CurrentWeather.Time)
	temperature := fmt.Sprintf("%v", weatherData.CurrentWeather.Temperature)
	windSpeed := fmt.Sprintf("%v", weatherData.CurrentWeather.WindSpeed)
	weatherCode := weatherData.CurrentWeather.WeatherCode
	direction := getWindDirection(int(weatherData.CurrentWeather.WindDirection))

	fmt.Printf("Lokalizacja: %s \n", city)
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

	data := []string{
		formattedTime,
		city,
		description,
		temperature,
		windSpeed,
		direction,
	}

	writeToFile("data.csv", data)

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
	return warsawTime.Format("2006-01-02 15:04 MST")
}

func main() {

	latitude := flag.String("latitude", "", "Latitude Value")
	longitude := flag.String("longitude", "", "Longitude Value")

	flag.Parse()

	if *latitude == "" || *longitude == "" {
		fmt.Println("Please provide latitude and longitude for location you want to check weather")
		fmt.Println("Example: <binary> --latitude=54.52 --longitude=18.53")
	}

	// Start a goroutine to fetch weather data every hour
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			weatherData, err := fetchWeather(*latitude, *longitude)
			if err != nil {
				fmt.Println("Error fetching weather data:", err)
			}
			city := getCity(*latitude, *longitude)

			handleWeatherData(weatherData, city)
			<-ticker.C
		}
	}()
	// Wait for a termination signal
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel
	fmt.Println("Received termination signal. Exiting...")
}
