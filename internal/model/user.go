package model

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	gorm.Model
	Name     string   `gorm:"not null" json:"name"`
	Email    string   `gorm:"uniqueIndex;not null" json:"email"`
	Password string   `gorm:"not null" json:"-"`
	Role     UserRole `gorm:"type:varchar(16);default:user;not null" json:"role"`
}

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
