package auth

import (
	"database/sql"
	"time"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository/user"
)

type CreateUserDTO struct {
	Name         string
	Username     string
	Email        string
	PasswordHash string
	Pro          int
	Admin        int
}

type AuthService struct {
	userRepo *user.UserRepository
}

func NewAuthService(userRepo *user.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) GetUserByUsername(username string) (*entity.Users, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *AuthService) GetUserByUserId(userId int) (*entity.Users, error) {
	return s.userRepo.GetByUserId(userId)
}

func (s *AuthService) GetAllUsers(page, pageSize int) ([]entity.Users, error) {
	return s.userRepo.GetAll(page, pageSize)
}

func (s *AuthService) CreateUser(dto CreateUserDTO) (*entity.Users, error) {
	user := &entity.Users{
		Username:        dto.Username,
		PasswordHash:    dto.PasswordHash,
		Email:           sql.NullString{String: dto.Email, Valid: true},
		Name:            sql.NullString{String: dto.Name, Valid: true},
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

func (s *AuthService) UpdateUserFromParams(UserId int, ActionType string, Value int) error {
	switch ActionType {
	case "admin":
		err := s.userRepo.UpdateUserAdmin(UserId, Value)
		if err != nil {
			return err
		}
	case "pro":
		err := s.userRepo.UpdateUserPro(UserId, Value)
		if err != nil {
			return err
		}
	case "block":
		err := s.userRepo.UpdateUserBlock(UserId, Value)
		if err != nil {
			return err
		}
	case "del":
		err := s.userRepo.UpdateUserDel(UserId, Value)
		if err != nil {
			return err
		}
	}
	return nil
}
