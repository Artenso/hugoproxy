package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Handler struct {
	GeoServ *GeoService
	Storage *Storage
}

func NewHandler(geoserv *GeoService, storage *Storage) *Handler {
	return &Handler{
		GeoServ: geoserv,
		Storage: storage,
	}
}

// @Summary      Log in
// @Description  returns JWT if you are registered user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input   body      RequestAuth  true  "registration data"
// @Success      200
// @Failure      400
// @Router       /api/login [post]
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	input := &RequestAuth{}

	data, _ := io.ReadAll(r.Body)
	json.Unmarshal(data, input)

	if input.Name == "" || input.Pass == "" {
		http.Error(w, "Missing username or password.", http.StatusBadRequest)
		return
	}

	if h.Storage.IsRegistered(input) {
		resp, err := json.Marshal(
			&ResponseLogin{
				Token: GenerateToken(input.Name),
			},
		)
		if err != nil {
			fmt.Println(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Register if you are not or check your username and password"))
}

// @Summary      Registration
// @Description  Saves your username and password in db
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input   body      RequestAuth  true  "registration data"
// @Success      200
// @Failure      400
// @Router       /api/register [post]
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	input := &RequestAuth{}

	data, _ := io.ReadAll(r.Body)
	json.Unmarshal(data, input)

	if input.Name == "" || input.Pass == "" {
		http.Error(w, "Missing username or password.", http.StatusBadRequest)
		return
	}

	user, err := NewUser(input.Name, input.Pass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	h.Storage.Add(user)

	w.WriteHeader(http.StatusOK)
}

// @Summary      GeoCode
// @Security     ApiKeyAuth
// @Description  Get full address info by coordinates
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        lat   query      string  true  "latitude"
// @Param        lng   query      string  true  "longitude"
// @Success      200  {object}  ResponseAddress
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Router       /api/address/geocode [post]
func (h *Handler) GeoCodeHandler(w http.ResponseWriter, r *http.Request) {
	requestAddressGeocode := RequestAddressGeocode{
		Lat: r.FormValue("lat"),
		Lng: r.FormValue("lng"),
	}

	if requestAddressGeocode.Lat == "" || requestAddressGeocode.Lng == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.GeoServ.GeoCode(requestAddressGeocode.Lat, requestAddressGeocode.Lng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// @Summary      Adress Search
// @Security     ApiKeyAuth
// @Description  Get full address info by its part
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        query   query      string  true  "part of address"
// @Success      200  {object}  ResponseAddress
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Router       /api/address/search [post]
func (h *Handler) AddressSearchHandler(w http.ResponseWriter, r *http.Request) {
	requestAddressSearch := RequestAddressSearch{
		Query: r.FormValue("query"),
	}

	if len(requestAddressSearch.Query) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.GeoServ.AddressSearch(requestAddressSearch.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
