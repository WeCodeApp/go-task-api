package models

import (
	"time"
)

type Client struct {
	ID        uint `gorm:"primaryKey"`
	Issuer    string
	Secret    string
	CreatedAt time.Time
}
