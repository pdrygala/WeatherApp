package city

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

func GetCity(latitude string, longitude string) string {
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
