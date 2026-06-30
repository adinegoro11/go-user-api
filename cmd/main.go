package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/adinegoro11/go-user-api/config"
	"github.com/adinegoro11/go-user-api/internal/handler"
	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
	"github.com/adinegoro11/go-user-api/internal/service"
	"github.com/adinegoro11/go-user-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Product{}); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	if err := db.AutoMigrate(
		&model.Customer{},
		&model.Subscription{},
		&model.SubscriptionProduct{},
		&model.Invoice{},
		&model.InvoiceItem{},
		&model.EmailLog{},
		&model.OTTRequest{},
		&model.DomainEvent{},
	); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)
	eventRepo := repository.NewEventRepository(db)
	emailLogRepo := repository.NewEmailLogRepository(db)
	ottRepo := repository.NewOTTRequestRepository(db)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "change-me-in-production"
	}

	authService := service.NewAuthService(userRepo, jwtSecret)
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	eventPublisher := service.NewDBEventPublisher(eventRepo)
	notificationService := service.NewEmailNotificationService(emailLogRepo)
	ottService := service.NewDefaultOTTService(ottRepo, nil)
	billingService := service.NewBillingService(invoiceRepo, customerRepo, subscriptionRepo, notificationService)
	customerService := service.NewCustomerService(customerRepo, eventPublisher)
	subscriptionService := service.NewSubscriptionService(customerRepo, subscriptionRepo, eventPublisher, notificationService, ottService, billingService)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	customerHandler := handler.NewCustomerHandler(customerService)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	billingHandler := handler.NewBillingHandler(billingService)

	router := gin.Default()
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	routes.SetupRoutes(router, authHandler, userHandler, productHandler, customerHandler, subscriptionHandler, billingHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
