package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/j94veron/auth-service-insu/internal/auth"
	"github.com/j94veron/auth-service-insu/internal/handlers"
	"github.com/j94veron/auth-service-insu/internal/middlewares"
	"github.com/j94veron/auth-service-insu/internal/models"
	"github.com/j94veron/auth-service-insu/internal/role"
	"github.com/j94veron/auth-service-insu/internal/user"
	"github.com/j94veron/auth-service-insu/pkg/redis"
	"github.com/j94veron/auth-service-insu/pkg/token"
)

func main() {
	// Cargar archivo .env desde /app (desde Docker)
	if err := godotenv.Load("/app/.env"); err != nil {
		log.Printf("Error loading .env file: %v", err)
	} else {
		log.Println("Successfully loaded .env from /app")
	}
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file")
	} else {
		log.Println("Successfully loaded .env from root")
	}

	// Conectar a la base de datos
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrar modelos
	db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})

	// Inicializar Redis (opcional)
	redisClient := redis.NewClient(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		0,
	)

	// Inicializar servicios y repositorios
	userRepo := user.NewRepository(db)
	roleRepo := role.NewRepository(db)

	tokenService := token.NewTokenService(
		os.Getenv("JWT_ACCESS_SECRET"),
		os.Getenv("JWT_REFRESH_SECRET"),
	)

	authService := auth.NewService(userRepo, tokenService, redisClient)

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo)
	roleHandler := handlers.NewRoleHandler(roleRepo)

	// Inicializar middlewares
	authMiddleware := middlewares.NewAuthMiddleware(tokenService, redisClient)
	permMiddleware := middlewares.NewPermissionMiddleware(userRepo, roleRepo)

	// Configurar router
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

	// Rutas públicas
	r.POST("/api/login", authHandler.Login)
	r.POST("/api/refresh", authHandler.Refresh)

	// Rutas protegidas
	api := r.Group("/api", authMiddleware.AuthRequired())
	{
		// Usuarios
		api.GET("/users", permMiddleware.HasPermission("/api/users"), userHandler.List)
		api.GET("/users/:id", permMiddleware.HasPermission("/api/users"), userHandler.GetByID)
		api.POST("/users", permMiddleware.HasPermission("/api/users"), userHandler.Create)
		api.PUT("/users/:id", permMiddleware.HasPermission("/api/users"), userHandler.Update)
		api.DELETE("/users/:id", permMiddleware.HasPermission("/api/users"), userHandler.Delete)

		// Roles
		api.GET("/roles", permMiddleware.HasPermission("/api/roles"), roleHandler.List)
		api.GET("/roles/:id", permMiddleware.HasPermission("/api/roles"), roleHandler.GetByID)
		api.POST("/roles", permMiddleware.HasPermission("/api/roles"), roleHandler.Create)
		api.PUT("/roles/:id", permMiddleware.HasPermission("/api/roles"), roleHandler.Update)
		api.DELETE("/roles/:id", permMiddleware.HasPermission("/api/roles"), roleHandler.Delete)
	}

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
