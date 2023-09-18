package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"pdrygala.com/weather/transport/rest/api"
)

func NewServer() http.Handler {
	r := chi.NewRouter()

	//Defines Routes
	r.Get("/latest", api.GetLatestWeather)
	r.Get("/bulk", api.GetAllData)
	return r
}
