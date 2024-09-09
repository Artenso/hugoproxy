package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

const (
	apiKey          = "93e5303e65aee396ec4015c293dd86b8240f185c"
	secretKey       = "f9a6d64b3de023a9fd458dd2b9594124fa346b93"
	swaggerTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <script src="//unpkg.com/swagger-ui-dist@3/swagger-ui-standalone-preset.js"></script>
    <!-- <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui-standalone-preset.js"></script> -->
    <script src="//unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
    <!-- <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui-bundle.js"></script> -->
    <link rel="stylesheet" href="//unpkg.com/swagger-ui-dist@3/swagger-ui.css" />
    <!-- <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui.css" /> -->
	<style>
		body {
			margin: 0;
		}
	</style>
    <title>Swagger</title>
</head>
<body>
    <div id="swagger-ui"></div>
    <script>
        window.onload = function() {
          SwaggerUIBundle({
            url: "/public/swagger.json?{{.Time}}",
            dom_id: '#swagger-ui',
            presets: [
              SwaggerUIBundle.presets.apis,
              SwaggerUIStandalonePreset
            ],
            layout: "StandaloneLayout"
          })
        }
    </script>
</body>
</html>
`
)

type RequestAddressSearch struct {
	Query string `json:"query"`
}

// swagger:model
type ResponseAddress struct {
	// Addresses
	//
	// required: true
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
		// swagger:operation POST /api/address/search AddressSearch
		//
		// Get full address info by its part
		//
		// ---
		// produces:
		// - application/json
		// parameters:
		// - name: query
		//   in: query
		//   description: part of address
		//   required: true
		//   type: string
		// responses:
		//   '200':
		//     description: full addresses
		//     schema:
		//       type: array
		//       items:
		//         "$ref": "#/definitions/ResponseAddress"
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
		// swagger:operation POST /api/address/geocode GeoCode
		//
		// Get full address info by coordinates
		//
		// ---
		// produces:
		// - application/json
		// parameters:
		// - name: lat
		//   in: query
		//   description: latitude
		//   required: true
		//   type: string
		// - name: lng
		//   in: query
		//   description: longitude
		//   required: true
		//   type: string
		// responses:
		//   '200':
		//     description: full addresses
		//     schema:
		//       type: array
		//       items:
		//         "$ref": "#/definitions/ResponseAddress"
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
	})

	//SwaggerUI
	r.Get("/swagger/index.html", swaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))).ServeHTTP(w, r)
	})
	http.ListenAndServe(":8080", r)
}

func swaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.New("swagger").Parse(swaggerTemplate)
	if err != nil {
		return
	}
	err = tmpl.Execute(w, struct {
		Time int64
	}{
		Time: time.Now().Unix(),
	})
	if err != nil {
		return
	}
}
