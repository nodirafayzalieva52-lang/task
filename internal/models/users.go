package models

import "time"

const (
	UserRole  = "user"
	AdminRole = "admin"
)

type User struct {
	ID        int
	Name      string
	Phone     string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
}

type Order struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
}
