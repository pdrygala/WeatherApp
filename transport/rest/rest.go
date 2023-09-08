package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"pdrygala.com/weather/transport/rest/api"
)

func NewServer() http.Handler {
	r := chi.NewRouter()

	//Defines Routes
	r.Get("/", api.GetLatestWeather)
	return r
}
