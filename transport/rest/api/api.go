package api

import (
	"net/http"
	"time"

	"pdrygala.com/weather/service/mysql"
)

func GetLatestWeather(w http.ResponseWriter, r *http.Request) {
	dns := "foobar:password@tcp(127.0.0.1:3306)/db"
	conn, err := mysql.NewMySQL(dns)
	if err != nil {
		panic(err.Error())
	}
	response := conn.GetLatestRecord()
	w.Write([]byte(response))
}

func GetResults(w http.ResponseWriter, r *http.Request) {
	//get queryparams from
	startTimeStr := r.URL.Query().Get("startTime")
	endTimeStr := r.URL.Query().Get("endTime")

	var response string

	//If params
	if startTimeStr != "" && endTimeStr != "" {
		//Parse startTime to time.Time
		startTime, err := time.Parse("2006-01-02T15:04:00Z", startTimeStr)
		if err != nil {
			panic(err.Error())
		}

		endTime, err := time.Parse("2006-01-02T15:04:00Z", endTimeStr)
		if err != nil {
			panic(err.Error())
		}

		response = mysql.GetResultsByTimeRange(startTime, endTime)
	} else {
		response = mysql.GetAllRecords()
	}

	w.Write([]byte(response))
}
