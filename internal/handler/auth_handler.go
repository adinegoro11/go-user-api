package handler

import (
	"log/slog"
	"net/http"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (dto.UserResponse, error)
	Login(req dto.LoginRequest) (dto.AuthResponse, error)
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("register request validation failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		switch err {
		case service.ErrInvalidRegisterInput:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case service.ErrEmailAlreadyRegistered:
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		default:
			slog.Error("register failed", "email", req.Email, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to register"})
		}
		return
	}

	slog.Info("register handler success", "user_id", user.ID, "email", user.Email)
	c.JSON(http.StatusCreated, gin.H{"message": "register success", "user": user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := h.authService.Login(req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login success", "data": result})
}
