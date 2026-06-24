package models

import "time"

const (
	UserRole  = "user"
	AdminRole = "admin"
)

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
}
