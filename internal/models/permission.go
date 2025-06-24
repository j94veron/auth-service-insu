package models

import "time"

type Permission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Resource    string    `json:"resource"` // Resource name (ej: "users")
	Endpoint    string    `json:"endpoint"` // specific endpoint (ej: "/api/users")
	Method      string    `json:"method"`   // Method HTTP (GET, POST, etc.)
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
