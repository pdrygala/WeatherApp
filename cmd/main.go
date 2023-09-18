package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pdrygala.com/weather/core/city"
	"pdrygala.com/weather/core/weather"
	"pdrygala.com/weather/transport/rest"
)

func main() {
	// TODO
	// Add GET with query params by date
	// Add Get with all data from table
	// TBD...

	// Create the REST API server
	apiHandler := rest.NewServer()

	// Start the server in a goroutine
	go func() {
		if err := http.ListenAndServe(":8080", apiHandler); err != nil {
			fmt.Println("HTTP server error:", err)
		}
	}()

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
			weatherData, err := weather.FetchWeather(*latitude, *longitude)
			if err != nil {
				fmt.Println("Error fetching weather data:", err)
			}
			city := city.GetCity(*latitude, *longitude)

			weather.HandleWeatherData(weatherData, city)
			<-ticker.C
		}
	}()
	// Wait for a termination signal
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel
	fmt.Println("Received termination signal. Exiting...")

}
