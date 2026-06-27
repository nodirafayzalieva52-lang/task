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

type UpdateRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func (r *UpdateRequest) Validate() error {
	if r.Name == "" || r.Phone == "" {
		return errors.New("name and phone are required")
	}
	return nil
}

type CreateOrderRequest struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

func (r *CreateOrderRequest) Validate() error {
	if r.Amount <= 0 || r.Description == "" {
		return errors.New("Validation error: amount and description are required")
	}
	return nil
}

type UpdateOrderRequest struct {
	Description string  `json:"description"`
	Amount      float64 `json:"ammount"`
	Status      string  `json:"status"`
}

func (r *UpdateOrderRequest) Validate() error {
	if r.Description == "" || r.Amount <= 0 || r.Status == "" {
		return errors.New("validation error: description, amount and status are required")
	}
	return nil
}

type UpdateRoleRequest struct {
	Role string `json:"role"`
}

func (r *UpdateRoleRequest) Validate() error {
	if r.Role != UserRole && r.Role != AdminRole {
		return errors.New("validation error: invalid role")
	}
	return nil
}

type UpdateUserPassword struct {
	OldPassword string `json:"old_password"` 
	NewPassword string `json:"new_password"` 
}

func (r *UpdateUserPassword) Validate() error {
	if r.OldPassword == "" || r.NewPassword == ""{
		return errors.New("Validation error: invalid password")
	}
	return nil
}