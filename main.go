package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	"invoice-api/internal/db"
	"invoice-api/internal/handler"
	"invoice-api/internal/middleware"
)

func main() {
	// Load .env only in development (tidak akan ada file .env di Render)
	if os.Getenv("RENDER") == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found (this is normal in production)")
		}
	}

	// Debug environment variables
	log.Printf("Environment check:")
	log.Printf("PORT: %s", os.Getenv("PORT"))
	log.Printf("DATABASE_URL exists: %v", os.Getenv("DATABASE_URL") != "")
	log.Printf("POSTGRES_DSN exists: %v", os.Getenv("POSTGRES_DSN") != "")
	log.Printf("JWT_SECRET exists: %v", os.Getenv("JWT_SECRET") != "")

	// Initialize database
	log.Println("Initializing database...")
	db.Init()
	log.Println("✅ Database initialized successfully")

	// Set Gin mode based on environment
	if os.Getenv("RENDER") != "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS configuration - Update for production
	corsConfig := cors.DefaultConfig()
	
	// Allow multiple origins for development and production
	allowedOrigins := []string{"http://localhost:3000"}
	
	// Add production frontend URL if exists
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}
	
	// In production, you might want to allow all origins (be careful!)
	if os.Getenv("RENDER") != "" {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = allowedOrigins
	}
	
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))

	// Health check endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from Go API!",
			"status":  "healthy",
			"env":     gin.Mode(),
		})
	})

	// Health check for monitoring
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"timestamp": "healthy",
		})
	})

	// Route public
	r.POST("/api/register", handler.Register)
	r.POST("/api/login", handler.Login)

	// Debug endpoint (only in development)
	if os.Getenv("RENDER") == "" {
		r.GET("/debug-token", func(c *gin.Context) {
			tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
			secret := os.Getenv("JWT_SECRET")

			fmt.Println("Raw token :", tokenStr)
			fmt.Println("JWT_SECRET:", secret)

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			valid := token != nil && token.Valid
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			c.JSON(200, gin.H{
				"valid": valid,
				"err":   errMsg,
			})
		})
	}

	// Route protected
	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired())

	{
		protected.POST("/invoices", handler.CreateInvoice)
		protected.GET("/invoices", handler.ListInvoices)
		protected.GET("/invoices/:id", handler.GetInvoice)
		protected.PUT("/invoices/:id", handler.UpdateInvoice)
		protected.DELETE("/invoices/:id", handler.DeleteInvoice)
	}

	// Port configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port untuk development
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📍 Environment: %s", gin.Mode())
	
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}