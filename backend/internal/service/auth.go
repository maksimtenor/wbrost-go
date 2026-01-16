package service

import (
	"database/sql"
	"time"
	"wbrost-go/internal/repository"
)

type CreateUserDTO struct {
	Name         string
	Username     string
	Email        string
	PasswordHash string
	Pro          int // Изменено с ProAccount string на Pro int
	Admin        int
}

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) GetUserByUsername(username string) (*repository.User, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *AuthService) CreateUser(dto CreateUserDTO) (*repository.User, error) {
	user := &repository.User{
		Username:        dto.Username,
		PasswordHash:    dto.PasswordHash,
		Email:           sql.NullString{String: dto.Email, Valid: true}, // Valid: true!
		Name:            sql.NullString{String: dto.Name, Valid: true},  // Valid: true!
		Pro:             dto.Pro,
		Admin:           dto.Admin,
		Taxes:           0,
		Block:           0,
		Phone:           sql.NullString{String: "", Valid: false},
		WbKey:           sql.NullString{String: "", Valid: false},
		OzonKey:         sql.NullString{String: "", Valid: false},
		U2782212Wbrosus: 0,
		OzonStatus:      0,
		Del:             0,
		LastLogin:       time.Now(),
	}

	return s.userRepo.Create(user)
}
