package main

import (
	"log"

	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/internal/auth"
	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/internal/database"
	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/internal/middleware"
	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/pkg/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(middleware.CORS(cfg.FrontendURL))

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg.JWTSecret, cfg.JWTExpirationHours)
	authHandler := auth.NewHandler(authService)

	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/me", middleware.AuthMiddleware(authService), authHandler.GetMe)
		authGroup.PATCH("/me", middleware.AuthMiddleware(authService), authHandler.UpdateMe)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := cfg.Port
	log.Printf("Server starting on port %s (auth slice)", port)
	log.Printf("Environment: %s", cfg.Environment)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
