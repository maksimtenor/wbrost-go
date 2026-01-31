package dto

import (
	"time"
)

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Name      *string   `json:"name"`
	Username  string    `json:"username"`
	Email     *string   `json:"email,omitempty"`
	Pro       int       `json:"pro"`
	WbKey     *string   `json:"wb_key,omitempty"`
	Taxes     int       `json:"taxes,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	Admin     int       `json:"admin"`
	CreatedAt time.Time `json:"created_at"`
}
