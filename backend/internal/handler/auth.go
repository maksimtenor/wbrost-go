package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"wbrost-go/internal/service"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtSecret   []byte
}

func NewAuthHandler(authService *service.AuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtSecret:   []byte(jwtSecret),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Валидация
	if req.Username == "" || req.Password == "" {
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Username and password are required"})
		return
	}

	// Получаем пользователя
	user, err := h.authService.GetUserByUsername(req.Username)
	if err != nil {
		// Возвращаем 400 вместо 401 для лучшей совместимости
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid username or password"})
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid username or password"})
		return
	}

	// ... остальной код без изменений
	// Генерируем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Формируем ответ
	response := LoginResponse{
		Token: tokenString,
		User: UserResponse{
			ID:        user.ID,
			Name:      getStringFromNullString(user.Name),
			Username:  user.Username,
			Email:     getStringFromNullString(user.Email),
			Pro:       user.Pro,
			Admin:     user.Admin,
			CreatedAt: user.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация
	validationErrors := make(map[string]string)
	if req.Name == "" {
		validationErrors["name"] = "Name is required"
	}
	if req.Username == "" {
		validationErrors["username"] = "Username is required"
	}
	if req.Email == "" {
		validationErrors["email"] = "Email is required"
	}
	if req.Password == "" || len(req.Password) < 6 {
		validationErrors["password"] = "Password must be at least 6 characters"
	}

	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ValidationErrors{Errors: validationErrors})
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Создаем пользователя - исправляем поле Pro
	user, err := h.authService.CreateUser(service.CreateUserDTO{
		Name:         req.Name,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Pro:          0, // По умолчанию не PRO (0 вместо "FREE")
		Admin:        0, // По умолчанию не админ
	})

	if err != nil {
		// Проверяем на уникальность
		if err.Error() == "username already exists" {
			validationErrors["username"] = "Username already taken"
		} else if err.Error() == "email already exists" {
			validationErrors["email"] = "Email already registered"
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		if len(validationErrors) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ValidationErrors{Errors: validationErrors})
			return
		}
	}

	// Генерируем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := LoginResponse{
		Token: tokenString,
		User: UserResponse{
			ID:        user.ID,
			Name:      getStringFromNullString(user.Name),
			Username:  user.Username,
			Email:     getStringFromNullString(user.Email),
			Pro:       user.Pro,
			Admin:     user.Admin,
			CreatedAt: user.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// Helper функция для конвертации sql.NullString в string
func getStringFromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// handler/auth.go
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Извлекаем токен из заголовка
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "No authorization header"})
		return
	}

	// Убираем "Bearer " префикс
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		respondWithJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Invalid token"})
		return
	}

	// Извлекаем данные из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Invalid token claims"})
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Invalid username in token"})
		return
	}

	// Получаем ПОСЛЕДНИЕ данные из БД
	user, err := h.authService.GetUserByUsername(username)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	// Формируем ответ с текущими данными
	response := UserResponse{
		ID:        user.ID,
		Name:      getStringFromNullString(user.Name),
		Username:  user.Username,
		Email:     getStringFromNullString(user.Email),
		Pro:       user.Pro,
		Admin:     user.Admin,
		CreatedAt: user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
