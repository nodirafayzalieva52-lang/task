package models

import "errors"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	if r.Name == "" || r.Email == "" || len(r.Password) < 8 {
		return errors.New("validation error: name and email are required")
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" || len(r.Password) < 8 {
		return errors.New("validation error: email and password are required")
	}
	return nil
}
