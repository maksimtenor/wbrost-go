package handler

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
	Name      string    `json:"name"` // string, не sql.NullString!
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"` // string!
	Pro       int       `json:"pro"`
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
