package dto

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRequest struct {
	UserId     int    `json:"userId"`
	ActionType string `json:"actionType"`
	Value      int    `json:"value"`
}
