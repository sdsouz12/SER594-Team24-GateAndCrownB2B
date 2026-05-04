package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/internal/aiproxy"
	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/internal/auth"
	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/internal/catalog"
	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/internal/database"
	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/internal/middleware"
	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/internal/orders"
	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/pkg/config"
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

	// Auth
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg.JWTSecret, cfg.JWTExpirationHours)
	authHandler := auth.NewHandler(authService)

	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.Register)
		authGroup.GET("/me", middleware.AuthMiddleware(authService), authHandler.GetMe)
		authGroup.PATCH("/me", middleware.AuthMiddleware(authService), authHandler.UpdateMe)
	}

	// Catalog
	catalogRepo := catalog.NewRepository(db)
	catalogService := catalog.NewService(catalogRepo, cfg.AIServiceURL)
	catalogHandler := catalog.NewHandler(catalogService)

	catalogGroup := router.Group("/api/catalog")
	{
		catalogGroup.GET("/products", catalogHandler.ListProducts)
		catalogGroup.POST("/search", catalogHandler.SemanticSearch)
	}

	// Orders
	ordersRepo := orders.NewRepository(db)
	ordersService := orders.NewService(ordersRepo)
	ordersHandler := orders.NewHandler(ordersService)

	ordersGroup := router.Group("/api/orders", middleware.AuthMiddleware(authService))
	{
		ordersGroup.POST("", ordersHandler.CreateOrder)
		ordersGroup.GET("", ordersHandler.GetMyOrders)
	}

	// AI proxy (forwards to Python AI-service)
	aiHandler := aiproxy.NewHandler(cfg.AIServiceURL)
	aiGroup := router.Group("/api/ai")
	{
		aiGroup.POST("/assist", aiHandler.Proxy("/assist"))
		aiGroup.POST("/rag", aiHandler.Proxy("/rag"))
		aiGroup.POST("/pipeline/ingest", aiHandler.Proxy("/pipeline/ingest"))
		aiGroup.GET("/pipeline/status", aiHandler.Proxy("/pipeline/status"))
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := cfg.Port
	log.Printf("Server starting on port %s", port)
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("AI Service URL: %s", cfg.AIServiceURL)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
