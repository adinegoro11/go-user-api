package handler

import (
	"net/http"
	"strconv"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/repository"
	"github.com/adinegoro11/go-user-api/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.userService.Me(userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role}})
}

func (h *UserHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := h.userService.Update(userID, req)
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update success", "user": gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role}})
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch users"})
		return
	}

	response := make([]gin.H, 0, len(users))
	for _, user := range users {
		response = append(response, gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role})
	}

	c.JSON(http.StatusOK, gin.H{"users": response})
}

func (h *UserHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	user, err := h.userService.FindByID(uint(id))
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role}})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	if err := h.userService.Delete(uint(id)); err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete success"})
}
