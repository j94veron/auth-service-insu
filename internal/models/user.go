package models

import "time"

type User struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Email          string    `json:"email" gorm:"unique"`
	Password       string    `json:"-" gorm:"not null"`
	UserName       string    `json:"userName"`
	Name           string    `json:"name"`
	LastName       string    `json:"lastName"`
	CommercialZone string    `json:"commercialZone"`
	Warehouse      string    `json:"warehouse"`
	RoleID         uint      `json:"roleId"`
	Role           Role      `json:"role"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	OtherWarehouse string    `json:"otherWarehouse"`
	Province       string    `json:"province"`
	Reports        string    `json:"reports"`
}
