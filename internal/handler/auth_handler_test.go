package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/service"
	"github.com/gin-gonic/gin"
)

type mockAuthService struct {
	registerFn func(req dto.RegisterRequest) (dto.UserResponse, error)
	loginFn    func(req dto.LoginRequest) (dto.AuthResponse, error)
}

func (m *mockAuthService) Register(req dto.RegisterRequest) (dto.UserResponse, error) {
	if m.registerFn != nil {
		return m.registerFn(req)
	}
	return dto.UserResponse{}, nil
}

func (m *mockAuthService) Login(req dto.LoginRequest) (dto.AuthResponse, error) {
	if m.loginFn != nil {
		return m.loginFn(req)
	}
	return dto.AuthResponse{}, nil
}

func setupRegisterRouter(svc AuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAuthHandler(svc)
	r.POST("/register", h.Register)
	return r
}

func TestAuthHandlerRegisterSuccess(t *testing.T) {
	r := setupRegisterRouter(&mockAuthService{
		registerFn: func(req dto.RegisterRequest) (dto.UserResponse, error) {
			return dto.UserResponse{ID: 1, Name: "John", Email: "john@example.com", Role: "user"}, nil
		},
	})

	body, _ := json.Marshal(gin.H{"name": "John", "email": "john@example.com", "password": "password123"})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestAuthHandlerRegisterDuplicateEmail(t *testing.T) {
	r := setupRegisterRouter(&mockAuthService{
		registerFn: func(req dto.RegisterRequest) (dto.UserResponse, error) {
			return dto.UserResponse{}, service.ErrEmailAlreadyRegistered
		},
	})

	body, _ := json.Marshal(gin.H{"name": "John", "email": "john@example.com", "password": "password123"})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestAuthHandlerRegisterServiceFailure(t *testing.T) {
	r := setupRegisterRouter(&mockAuthService{
		registerFn: func(req dto.RegisterRequest) (dto.UserResponse, error) {
			return dto.UserResponse{}, errors.New("db error")
		},
	})

	body, _ := json.Marshal(gin.H{"name": "John", "email": "john@example.com", "password": "password123"})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestAuthHandlerRegisterInvalidPayload(t *testing.T) {
	r := setupRegisterRouter(&mockAuthService{})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
