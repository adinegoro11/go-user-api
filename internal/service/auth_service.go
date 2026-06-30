package service

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrEmailAlreadyRegistered = errors.New("email already registered")
var ErrInvalidRegisterInput = errors.New("invalid register input")

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(req dto.RegisterRequest) (dto.UserResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if name == "" || email == "" {
		slog.Warn("register rejected due to invalid input", "email", email)
		return dto.UserResponse{}, ErrInvalidRegisterInput
	}

	existingUser, err := s.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		slog.Warn("register rejected due to duplicate email", "email", email)
		return dto.UserResponse{}, ErrEmailAlreadyRegistered
	}
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		slog.Error("register failed while checking existing user", "email", email, "error", err)
		return dto.UserResponse{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("register failed while hashing password", "email", email, "error", err)
		return dto.UserResponse{}, err
	}

	user := model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     model.RoleUser,
	}

	if err := s.userRepo.Create(&user); err != nil {
		slog.Error("register failed while creating user", "email", email, "error", err)
		return dto.UserResponse{}, err
	}

	slog.Info("register success", "user_id", user.ID, "email", user.Email)

	return toUserResponse(&user), nil
}

func (s *AuthService) Login(req dto.LoginRequest) (dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		return dto.AuthResponse{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return dto.AuthResponse{}, ErrInvalidCredentials
	}

	claims := model.JWTClaims{
		UserID: user.ID,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Token: signedToken,
		User:  toUserResponse(user),
	}, nil
}

func toUserResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}
}
