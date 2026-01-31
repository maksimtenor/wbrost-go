package dto

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationErrors struct {
	Errors map[string]string `json:"errors"`
}
