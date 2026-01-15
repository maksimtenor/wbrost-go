package entity

import (
	"fmt"
)

// LoginForm - точный аналог Yii2 LoginForm
type LoginForm struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`

	// Приватное поле для кэширования
	user *User
}

// Validate - аналог rules() и validatePassword()
func (f *LoginForm) Validate() error {
	// Правило 1: оба поля обязательны
	if f.Username == "" || f.Password == "" {
		return fmt.Errorf("username and password are required")
	}

	// Правило 2: rememberMe должен быть boolean
	// В Go это гарантировано типом bool

	// Правило 3: валидация пароля
	user, err := f.GetUser()
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil || !user.ValidatePassword(f.Password) {
		return fmt.Errorf("incorrect username or password")
	}

	return nil
}

// GetUser - аналог getUser() с кэшированием
func (f *LoginForm) GetUser() (*User, error) {
	if f.user == nil {
		// Здесь будет вызов репозитория
		// Пока заглушка
		f.user = &User{
			Username: f.Username,
			// Реальный поиск будет в репозитории
		}
	}
	return f.user, nil
}

// Login - аналог login()
func (f *LoginForm) Login() (bool, error) {
	if err := f.Validate(); err != nil {
		return false, err
	}

	user, err := f.GetUser()
	if err != nil {
		return false, err
	}

	// Здесь будет логика создания сессии
	// Пока просто возвращаем успех
	return true, nil
}

// AttributeLabels - аналог attributeLabels()
func (f *LoginForm) AttributeLabels() map[string]string {
	return map[string]string{
		"username":   "Логин",
		"password":   "Пароль",
		"rememberMe": "Запомнить меня",
	}
}
