package entity

import (
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Name      *string   `json:"name"` // string, не sql.NullString!
	Username  string    `json:"username"`
	Email     *string   `json:"email,omitempty"` // string!
	Pro       int       `json:"pro"`
	WbKey     *string   `json:"wb_key,omitempty"`
	Taxes     int       `json:"taxes,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	Admin     int       `json:"admin"`
	CreatedAt time.Time `json:"created_at"`
}

type SignupRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationErrors struct {
	Errors map[string]string `json:"errors"`
}
