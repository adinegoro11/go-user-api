package service

import (
	"errors"
	"testing"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockUserRepository struct {
	createFn      func(user *model.User) error
	findByEmailFn func(email string) (*model.User, error)
}

func (m *mockUserRepository) Create(user *model.User) error {
	if m.createFn != nil {
		return m.createFn(user)
	}
	return nil
}

func (m *mockUserRepository) FindByEmail(email string) (*model.User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(email)
	}
	return nil, repository.ErrUserNotFound
}

func (m *mockUserRepository) FindByID(id uint) (*model.User, error) {
	return nil, repository.ErrUserNotFound
}

func (m *mockUserRepository) FindAll() ([]model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) Update(user *model.User) error {
	return nil
}

func (m *mockUserRepository) Delete(id uint) error {
	return nil
}

func TestAuthServiceRegisterSuccess(t *testing.T) {
	var createdUser model.User
	repo := &mockUserRepository{
		findByEmailFn: func(email string) (*model.User, error) {
			if email != "john@example.com" {
				t.Fatalf("unexpected normalized email: %s", email)
			}
			return nil, repository.ErrUserNotFound
		},
		createFn: func(user *model.User) error {
			createdUser = *user
			createdUser.ID = 1
			user.ID = 1
			return nil
		},
	}

	svc := NewAuthService(repo, "secret")
	resp, err := svc.Register(dto.RegisterRequest{
		Name:     "  John Doe  ",
		Email:    "  JOHN@EXAMPLE.COM  ",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.ID != 1 {
		t.Fatalf("expected ID 1, got %d", resp.ID)
	}
	if createdUser.Name != "John Doe" {
		t.Fatalf("expected trimmed name, got %q", createdUser.Name)
	}
	if createdUser.Email != "john@example.com" {
		t.Fatalf("expected normalized email, got %q", createdUser.Email)
	}
	if createdUser.Role != model.RoleUser {
		t.Fatalf("expected role user, got %q", createdUser.Role)
	}
	if createdUser.Password == "password123" {
		t.Fatalf("expected hashed password, got plain text")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte("password123")); err != nil {
		t.Fatalf("expected valid bcrypt hash, got %v", err)
	}
}

func TestAuthServiceRegisterDuplicateEmail(t *testing.T) {
	createCalled := false
	repo := &mockUserRepository{
		findByEmailFn: func(email string) (*model.User, error) {
			return &model.User{Model: gorm.Model{ID: 22}, Email: email}, nil
		},
		createFn: func(user *model.User) error {
			createCalled = true
			return nil
		},
	}

	svc := NewAuthService(repo, "secret")
	_, err := svc.Register(dto.RegisterRequest{
		Name:     "John",
		Email:    "john@example.com",
		Password: "password123",
	})
	if !errors.Is(err, ErrEmailAlreadyRegistered) {
		t.Fatalf("expected ErrEmailAlreadyRegistered, got %v", err)
	}
	if createCalled {
		t.Fatalf("expected Create not to be called when email exists")
	}
}

func TestAuthServiceRegisterInvalidInput(t *testing.T) {
	repo := &mockUserRepository{}
	svc := NewAuthService(repo, "secret")

	_, err := svc.Register(dto.RegisterRequest{
		Name:     "   ",
		Email:    "   ",
		Password: "password123",
	})
	if !errors.Is(err, ErrInvalidRegisterInput) {
		t.Fatalf("expected ErrInvalidRegisterInput, got %v", err)
	}
}

func TestAuthServiceRegisterFindByEmailFailure(t *testing.T) {
	repo := &mockUserRepository{
		findByEmailFn: func(email string) (*model.User, error) {
			return nil, errors.New("database unavailable")
		},
	}
	svc := NewAuthService(repo, "secret")

	_, err := svc.Register(dto.RegisterRequest{
		Name:     "John",
		Email:    "john@example.com",
		Password: "password123",
	})
	if err == nil || err.Error() != "database unavailable" {
		t.Fatalf("expected database error, got %v", err)
	}
}
