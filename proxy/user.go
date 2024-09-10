package main

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name string `json:"username"`
	Pass []byte `json:"password"`
}

func NewUser(name, pass string) (*User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(pass), 16)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name: name,
		Pass: password,
	}

	return user, nil
}
