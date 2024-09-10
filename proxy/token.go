package main

import "github.com/go-chi/jwtauth"

const (
	Secret = "SamvelTheBest))"
)

var tokenAuth *jwtauth.JWTAuth

func GenerateToken(name string) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"username": name})
	return tokenString
}
