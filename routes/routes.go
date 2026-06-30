package routes

import (
	"github.com/adinegoro11/go-user-api/internal/handler"
	"github.com/adinegoro11/go-user-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, productHandler *handler.ProductHandler) {
	api := router.Group("/api")
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)
	api.GET("/products", productHandler.FindAll)
	api.GET("/products/:id", productHandler.FindByID)

	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware())
	auth.GET("/me", userHandler.Me)
	auth.PUT("/me", userHandler.Update)

	adminProducts := auth.Group("/products")
	adminProducts.Use(middleware.RoleMiddleware("admin"))
	adminProducts.POST("", productHandler.Create)
	adminProducts.PUT("/:id", productHandler.Update)
	adminProducts.DELETE("/:id", productHandler.Delete)

	admin := auth.Group("/admin")
	admin.Use(middleware.RoleMiddleware("admin"))
	admin.GET("/users", userHandler.GetAll)
	admin.GET("/users/:id", userHandler.Detail)
	admin.DELETE("/users/:id", userHandler.Delete)
}
