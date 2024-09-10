package main

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

type ResponseLogin struct {
	Token string `json:"token"`
}

type RequestAuth struct {
	Name string `json:"username"`
	Pass string `json:"password"`
}
