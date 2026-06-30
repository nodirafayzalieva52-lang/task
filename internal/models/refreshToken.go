package models

import "time"

type RefreshToken struct {
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
}
