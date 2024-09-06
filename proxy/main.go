package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	apiKey    = ""
	secretKey = ""
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

func main() {
	geoserv := NewGeoService(apiKey, secretKey)
	r := chi.NewRouter()

	r.Post("/api/address/search", func(w http.ResponseWriter, r *http.Request) {
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
	})

	r.Post("/api/address/geocode", func(w http.ResponseWriter, r *http.Request) {
		requestAddressGeocode := RequestAddressGeocode{
			Lat: r.FormValue("lat"),
			Lng: r.FormValue("lng"),
		}

		fmt.Println(requestAddressGeocode)

		if requestAddressGeocode.Lat == "" || requestAddressGeocode.Lng == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := geoserv.GeoCode(requestAddressGeocode.Lat, requestAddressGeocode.Lng)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(res)

		data, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)

		// for _, addres := range res {
		// 	w.Write([]byte(fmt.Sprintf(
		// 		"Координаты %s %s указывают на город: %s",
		// 		addres.Lat,
		// 		addres.Lon,
		// 		addres.City,
		// 	)),
		// 	)
		// }
	})

	http.ListenAndServe(":8080", r)
}

// const content = ``

// func WorkerTest() {
// 	t := time.NewTicker(1 * time.Second)
// 	var b byte = 0
// 	for {
// 		select {
// 		case <-t.C:
// 			err := os.WriteFile("/app/static/_index.md", []byte(fmt.Sprintf(content, b)), 0644)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			b++
// 		}
// 	}
// }
