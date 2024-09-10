package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/httputil"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware

	_ "test/docs"
)

const (
	apiKey    = "93e5303e65aee396ec4015c293dd86b8240f185c"
	secretKey = "f9a6d64b3de023a9fd458dd2b9594124fa346b93"
)

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type ResponseAddress struct {
	Addresses []*Address `json:"addresses"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

// @title Proxy
// @version 1.0
// @description Documentation of Proxy API.

// @host localhost:8080
// @BasePath /
func main() {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	r.Post("/api/address/search", AddressSearchHandler)

	r.Post("/api/address/geocode", GeoCodeHandler)

	http.ListenAndServe(":8080", r)
}

// @Summary      GeoCode
// @Description  Get full address info by coordinates
// @Tags         geocode
// @Accept       json
// @Produce      json
// @Param        lat   query      string  true  "latitude"
// @Param        lng   query      string  true  "longitude"
// @Success      200  {object}  ResponseAddress
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       //api/address/geocode [post]
func GeoCodeHandler(w http.ResponseWriter, r *http.Request) {
	geoserv := NewGeoService(apiKey, secretKey)

	requestAddressGeocode := RequestAddressGeocode{
		Lat: r.FormValue("lat"),
		Lng: r.FormValue("lng"),
	}

	if requestAddressGeocode.Lat == "" || requestAddressGeocode.Lng == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := geoserv.GeoCode(requestAddressGeocode.Lat, requestAddressGeocode.Lng)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// @Summary      Adress Search
// @Description  Get full address info by its part
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        query   query      string  true  "part of address"
// @Success      200  {object}  ResponseAddress
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       //api/address/search [post]
func AddressSearchHandler(w http.ResponseWriter, r *http.Request) {
	geoserv := NewGeoService(apiKey, secretKey)

	requestAddressSearch := RequestAddressSearch{
		Query: r.FormValue("query"),
	}

	if len(requestAddressSearch.Query) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := geoserv.AddressSearch(requestAddressSearch.Query)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
