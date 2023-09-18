package api

import (
	"net/http"

	"pdrygala.com/weather/service/mysql"
)

func GetLatestWeather(w http.ResponseWriter, r *http.Request) {
	response := mysql.GetLatestRecord()
	w.Write([]byte(response))
}

func GetAllData(w http.ResponseWriter, r *http.Request) {
	respone := mysql.GetAllRecords()
	w.Write([]byte(respone))
}
