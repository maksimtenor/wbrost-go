package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository"
	"wbrost-go/internal/service"
	"wbrost-go/internal/wbapi"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService *service.AuthService
	userRepo    *repository.UserRepository
	jwtSecret   []byte
}

func NewAuthHandler(authService *service.AuthService, userRepo *repository.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
		jwtSecret:   []byte(jwtSecret),
	}
}

func (h *AuthHandler) GetUsersList(w http.ResponseWriter, r *http.Request) {
	// Получить пользователя по токену
	userCheck, err := h.getUserFromRequest(r)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Unauthorized"})
		return
	}
	if userCheck.Admin < 1 {
		respondWithJSON(w, http.StatusForbidden, entity.ErrorResponse{Error: "Вам сюда нельзя!"})
		return
	}
	// Параметры запроса
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10 // По умолчанию
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	var users []entity.User
	var totalCount int

	// Получение всех пользователей с пагинацией
	users, err = h.userRepo.GetAll(page, pageSize)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, entity.ErrorResponse{
			Error: "Failed to get articles: " + err.Error(),
		})
		return
	}

	// Получить общее количество
	totalCount, err = h.userRepo.GetCountAllUsers()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, entity.ErrorResponse{
			Error: "Failed to get total count: " + err.Error(),
		})
		return
	}

	// Форматировать ответ
	response := make([]map[string]interface{}, len(users))
	for i, user := range users {

		// Форматирование дат
		createdDate := user.CreatedAt.Format("2006-01-02")
		updatedDate := user.UpdatedAt.Format("2006-01-02")
		// Извлекаем значения из NullString
		phone := ""
		if user.Phone.Valid {
			phone = user.Phone.String
		}

		response[i] = map[string]interface{}{
			"id":          user.ID,
			"taxes":       user.Taxes,
			"username":    user.Username,
			"password":    user.Password,
			"email":       user.Email,
			"admin":       user.Admin,
			"block":       user.Block,
			"pro":         user.Pro,
			"name":        user.Name,
			"phone":       phone,
			"wb_key":      user.WbKey,
			"ozon_key":    user.OzonKey,
			"ozon_status": user.OzonStatus,
			"created_at":  createdDate,
			"updated_at":  updatedDate,
			"del":         user.Del,
			"last_login":  user.LastLogin,
		}
	}

	// Полный ответ с пагинацией
	fullResponse := map[string]interface{}{
		"data": response,
		"pagination": map[string]interface{}{
			"current_page": page,
			"page_size":    pageSize,
			"total_items":  totalCount,
			"total_pages":  (totalCount + pageSize - 1) / pageSize,
		},
	}

	respondWithJSON(w, http.StatusOK, fullResponse)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req entity.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Валидация
	if req.Username == "" || req.Password == "" {
		respondWithJSON(w, http.StatusBadRequest, entity.ErrorResponse{Error: "Username and password are required"})
		return
	}

	// Получаем пользователя
	user, err := h.authService.GetUserByUsername(req.Username)
	if err != nil {
		// Возвращаем 400 вместо 401 для лучшей совместимости
		respondWithJSON(w, http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid username or password"})
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		respondWithJSON(w, http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid username or password"})
		return
	}

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
	response := entity.LoginResponse{
		Token: tokenString,
		User: entity.UserResponse{
			ID:        user.ID,
			Name:      getStringPtrFromNullString(user.Name),
			Username:  user.Username,
			Email:     getStringPtrFromNullString(user.Email),
			Pro:       user.Pro,
			Taxes:     user.Taxes,
			WbKey:     getStringPtrFromNullString(user.WbKey),
			Phone:     getStringPtrFromNullString(user.Phone),
			Admin:     user.Admin,
			CreatedAt: user.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req entity.SignupRequest
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
		json.NewEncoder(w).Encode(entity.ValidationErrors{Errors: validationErrors})
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
		Pro:          0, // По умолчанию не активирован Pro
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
			json.NewEncoder(w).Encode(entity.ValidationErrors{Errors: validationErrors})
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

	response := entity.LoginResponse{
		Token: tokenString,
		User: entity.UserResponse{
			ID:        user.ID,
			Name:      getStringPtrFromNullString(user.Name),
			Username:  user.Username,
			Email:     getStringPtrFromNullString(user.Email),
			Pro:       user.Pro,
			Taxes:     user.Taxes,
			WbKey:     getStringPtrFromNullString(user.WbKey),
			Phone:     getStringPtrFromNullString(user.Phone),
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
	json.NewEncoder(w).Encode(entity.ErrorResponse{Error: message})
}

// Helper функция для конвертации sql.NullString в string
func getStringPtrFromNullString(ns sql.NullString) *string {
	if ns.Valid {
		s := ns.String
		return &s
	}
	return nil
}

// GetCurrentUser - получаем текущего юзера
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Извлекаем токен из заголовка
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "No authorization header"})
		return
	}

	// Убираем "Bearer " префикс
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Invalid token"})
		return
	}

	// Извлекаем данные из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Invalid token claims"})
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Invalid username in token"})
		return
	}

	// Получаем ПОСЛЕДНИЕ данные из БД
	user, err := h.authService.GetUserByUsername(username)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, entity.ErrorResponse{Error: "User not found"})
		return
	}

	// Формируем ответ с текущими данными
	response := entity.UserResponse{
		ID:        user.ID,
		Name:      getStringPtrFromNullString(user.Name),
		Username:  user.Username,
		Email:     getStringPtrFromNullString(user.Email),
		Pro:       user.Pro,
		Taxes:     user.Taxes,
		WbKey:     getStringPtrFromNullString(user.WbKey),
		Phone:     getStringPtrFromNullString(user.Phone),
		Admin:     user.Admin,
		CreatedAt: user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetApiKeysStatus(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем токен из заголовка
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "No authorization header"})
		return
	}

	// 2. Извлекаем username из JWT (используем существующую логику как в GetCurrentUser)
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Invalid token"})
		return
	}

	// 3. Извлекаем username из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Invalid token claims"})
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Invalid username in token"})
		return
	}

	// 4. Получаем пользователя через authService (не userRepo!)
	user, err := h.authService.GetUserByUsername(username)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, entity.ErrorResponse{Error: "User not found"})
		return
	}

	// 5. Проверяем WB ключ
	wbStatus := map[string]interface{}{
		"has_token": user.WbKey.Valid && user.WbKey.String != "",
		"active":    false,
		"message":   "Не настроен",
	}

	if user.WbKey.Valid && user.WbKey.String != "" {
		// СОЗДАЕМ WB КЛИЕНТ И ПРОВЕРЯЕМ ТОКЕН
		wbClient := wbapi.NewWBClient(user.WbKey.String)
		isValid, err := wbClient.CheckToken()

		if err != nil {
			// Ошибка при проверке (сеть, timeout и т.д.)
			wbStatus["active"] = false
			wbStatus["message"] = "Ошибка проверки: " + err.Error()
		} else if isValid {
			// Токен рабочий!
			wbStatus["active"] = true
			wbStatus["message"] = "Активен"
		} else {
			// Токен невалидный (WB API вернул ошибку)
			wbStatus["active"] = false
			wbStatus["message"] = "Токен недействителен"
		}
		//// Здесь будет проверка через WB API
		//// Пока заглушка - всегда true
		//wbStatus["active"] = true
		//wbStatus["message"] = "Активен"
	}

	// 6. Возвращаем ответ
	response := map[string]interface{}{
		"wildberries": wbStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// 1. Проверяем авторизацию
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "No authorization header"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Invalid token"})
		return
	}

	// 2. Получаем username из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Invalid token claims"})
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Invalid username in token"})
		return
	}

	// 3. Получаем текущего пользователя
	currentUser, err := h.authService.GetUserByUsername(username)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, entity.ErrorResponse{Error: "User not found"})
		return
	}

	// 4. Парсим данные из запроса
	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		respondWithJSON(w, http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid request body"})
		return
	}

	// 5. Обновляем поля (только разрешенные)
	if name, ok := updateData["name"].(string); ok && name != "" {
		currentUser.Name = sql.NullString{String: name, Valid: true}
	}

	if email, ok := updateData["email"].(string); ok && email != "" {
		currentUser.Email = sql.NullString{String: email, Valid: true}
	}

	if phone, ok := updateData["phone"].(string); ok {
		currentUser.Phone = sql.NullString{String: phone, Valid: phone != ""}
	}

	if taxesStr, ok := updateData["taxes"].(float64); ok {
		currentUser.Taxes = int(taxesStr)
	}
	// 6. Обновляем WB токен (если изменился)
	if wbKey, ok := updateData["wb_key"].(string); ok {
		currentUser.WbKey = sql.NullString{String: wbKey, Valid: wbKey != ""}
	}

	// 7. Обновляем пароль (если указан новый)
	if password, ok := updateData["password"].(string); ok && password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to hash password"})
			return
		}
		currentUser.PasswordHash = string(hashedPassword)
	}

	// Замените пункт 8 в методе UpdateProfile:
	// 8. Сохраняем в БД
	err = h.userRepo.UpdateUser(currentUser)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to save user data: " + err.Error()})
		return
	}

	// 9. Возвращаем успешный ответ
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Profile updated successfully",
	})
}

func (h *AuthHandler) UpdateUserParams(w http.ResponseWriter, r *http.Request) {
	// 1. Проверяем авторизацию
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "No authorization header"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Invalid token"})
		return
	}

	// 2. Получаем username из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Invalid token claims"})
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Invalid username in token"})
		return
	}

	// 3. Получаем текущего пользователя
	currentUser, err := h.authService.GetUserByUsername(username)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, entity.ErrorResponse{Error: "User not found"})
		return
	}

	// 4. Проверяем права администратора
	if currentUser.Admin < 1 {
		respondWithJSON(w, http.StatusForbidden, entity.ErrorResponse{Error: "Admin rights required"})
		return
	}

	// Парсим JSON из тела запроса
	var req entity.UserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Проверяем обязательные поля
	if req.UserId == 0 || req.ActionType == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Обрабатываем данные
	err = h.authService.UpdateUserFromParams(req.UserId, req.ActionType, req.Value)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// 8. Возвращаем успешный ответ
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "User updated successfully",
	})
}

func (h *AuthHandler) getUserFromRequest(r *http.Request) (*repository.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("no authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username in token")
	}

	return h.authService.GetUserByUsername(username)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
