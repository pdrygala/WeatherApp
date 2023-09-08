package weather

type WeatherData struct {
	FormattedTime string
	City          string
	WeatherCode   int
	Description   string
	Temperature   string
	WindSpeed     string
	Direction     string
}

type WeatherDB struct {
	Id            int
	FormattedTime string
	City          string
	Description   string
	Temperature   string
	WindSpeed     string
	Direction     string
}
