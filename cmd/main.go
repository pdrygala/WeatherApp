package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pdrygala.com/weather/core/city"
	"pdrygala.com/weather/core/weather"
)

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
