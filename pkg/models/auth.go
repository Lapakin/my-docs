package models

import (
	"errors"
	"time"
)

type RegisterRequest struct {
	Username  string `json:"username" form:"username" db:"username"`
	Email     string `json:"email" form:"email" db:"email"`
	Password  string `json:"password" form:"password" db:"password"`
	FirstName string `json:"first_name" form:"first_name" db:"first_name"`
	LastName  string `json:"last_name" form:"last_name" db:"last_name"`
}

func (r *RegisterRequest) Validate() error {
	if r.Username == "" || r.Email == "" || r.Password == "" {
		return errors.New("username, email, and password are required")
	}
	return nil
}

type LoginRequest struct {
	Login    string `json:"login" form:"login"`
	Password string `json:"password" form:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Login == "" || r.Password == "" {
		return errors.New("login and password are required")
	}
	return nil
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" db:"old_password"`
	NewPassword string `json:"new_password" db:"new_password"`
}

func (r *ChangePasswordRequest) Validate() error {
	if r.OldPassword == "" || r.NewPassword == "" {
		return errors.New("old password and new password are required")
	}
	return nil
}

type RefreshTokenData struct {
	UserID    uint64
	Token     string
	ExpiresAt int64
	CreatedAt time.Time
}

type PasswordData struct {
	UserID     uint64
	Hash       string
	ModifiedAt time.Time
}
