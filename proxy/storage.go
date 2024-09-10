package main

import "golang.org/x/crypto/bcrypt"

type Storage struct {
	data []*User
}

func NewStorage() *Storage {
	return &Storage{
		data: make([]*User, 0, 16),
	}
}

func (s *Storage) Add(user *User) {
	s.data = append(s.data, user)
}

func (s *Storage) IsRegistered(input *RequestAuth) bool {
	for _, dbuser := range s.data {
		if dbuser.Name == input.Name {
			if bcrypt.CompareHashAndPassword(dbuser.Pass, []byte(input.Pass)) == nil {
				return true
			}
		}
	}

	return false
}
