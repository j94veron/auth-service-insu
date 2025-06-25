package main

import (
	"github.com/j94veron/auth-service-insu/internal/config"
	"github.com/j94veron/auth-service-insu/logger"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/j94veron/auth-service-insu/internal/auth"
	"github.com/j94veron/auth-service-insu/internal/handlers"
	"github.com/j94veron/auth-service-insu/internal/middlewares"
	"github.com/j94veron/auth-service-insu/internal/models"
	"github.com/j94veron/auth-service-insu/internal/role"
	"github.com/j94veron/auth-service-insu/internal/user"
	"github.com/j94veron/auth-service-insu/pkg/redis"
	"github.com/j94veron/auth-service-insu/pkg/token"
	"github.com/joho/godotenv"
)

func main() {
	// Upload .env files
	if err := godotenv.Load("../.env"); err != nil {
		logger.Logger.Error("Error loading .env file: " + err.Error())
	} else {
		logger.Logger.Info("Successfully loaded .env from /app")
	}
	if err := godotenv.Load(); err != nil {
		logger.Logger.Error("Error loading .env file: " + err.Error())
	} else {
		logger.Logger.Info("Successfully loaded .env from root")
	}

	//Connect to the database using the new package
	db, err := config.ConnectDB()
	if err != nil {
		logger.Logger.Error("Error connecting to database: " + err.Error())
	}

	// Auto-migrate models
	db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})

	//Initialize Redis (optional)
	redisClient := redis.NewClient(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		0,
	)

	// Initialize services and repositories
	userRepo := user.NewRepository(db)
	roleRepo := role.NewRepository(db)

	tokenService := token.NewTokenService(
		os.Getenv("JWT_ACCESS_SECRET"),
		os.Getenv("JWT_REFRESH_SECRET"),
	)

	authService := auth.NewService(userRepo, tokenService, redisClient)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo)
	roleHandler := handlers.NewRoleHandler(roleRepo)

	// Initialize middlewares
	authMiddleware := middlewares.NewAuthMiddleware(tokenService, redisClient)
	permMiddleware := middlewares.NewPermissionMiddleware(userRepo, roleRepo)

	// Configure router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Public routes
	r.POST("/api/login", authHandler.Login)
	r.POST("/api/refresh_token", authHandler.Refresh)

	// Protected routes
	api := r.Group("/api", authMiddleware.AuthRequired())
	{
		// User
		api.GET("/users", permMiddleware.HasPermission("/api/users"), userHandler.List)
		api.GET("/users/:id", permMiddleware.HasPermission("/api/users"), userHandler.GetByID)
		api.POST("/users", permMiddleware.HasPermission("/api/users"), userHandler.Create)
		api.PUT("/users/:id", permMiddleware.HasPermission("/api/users"), userHandler.Update)
		api.DELETE("/users/:id", permMiddleware.HasPermission("/api/users"), userHandler.Delete)

		// Role
		api.GET("/roles", permMiddleware.HasPermission("/api/roles"), roleHandler.List)
		api.GET("/roles/:id", permMiddleware.HasPermission("/api/roles"), roleHandler.GetByID)
		api.POST("/roles", permMiddleware.HasPermission("/api/roles"), roleHandler.Create)
		api.PUT("/roles/:id", permMiddleware.HasPermission("/api/roles"), roleHandler.Update)
		api.DELETE("/roles/:id", permMiddleware.HasPermission("/api/roles"), roleHandler.Delete)
	}

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	logger.Logger.Info("Starting server on " + port)
	if err := r.Run(port); err != nil {
		logger.Logger.Error("Server error: " + err.Error())
		log.Fatal(err)
	}
}
