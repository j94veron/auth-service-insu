package models

import "time"

type Permission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Resource    string    `json:"resource"` // Nombre del recurso (ej: "users")
	Endpoint    string    `json:"endpoint"` // Endpoint específico (ej: "/api/users")
	Method      string    `json:"method"`   // Método HTTP (GET, POST, etc.)
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
