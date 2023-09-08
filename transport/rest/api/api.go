package api

import (
	"net/http"

	"pdrygala.com/weather/service/mysql"
)

func GetLatestWeather(w http.ResponseWriter, r *http.Request) {
	response := mysql.GetLatestRecord()
	w.Write([]byte(response))
}
