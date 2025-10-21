package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         uint64     `json:"id" db:"id"`
	Email      string     `json:"email" db:"email"`
	Username   string     `json:"username" db:"username"`
	FirstName  string     `json:"first_name" db:"first_name"`
	LastName   string     `json:"last_name" db:"last_name"`
	Role       UserRole   `json:"role" db:"role"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	IsDeleted  bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time `json:"modified_at" db:"modified_at"`
}

type UserPassword struct {
	UserID       uint64     `json:"user_id" db:"user_id"`
	PasswordHash string     `json:"-" db:"password_hash"`
	ModifiedAt   *time.Time `json:"modified_at" db:"modified_at"`
}

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
	RoleGuest UserRole = "guest"
)

func NewUser(email, username, firstName, lastName string) *User {
	return &User{
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Role:      RoleUser,
		IsActive:  true,
		IsDeleted: false,
		CreatedAt: time.Now(),
	}
}

func (u *User) SetPassword(password string) (*UserPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &UserPassword{
		UserID:       u.ID,
		PasswordHash: string(hash),
	}, nil
}

func CheckPassword(passwordHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

type Users []*User
