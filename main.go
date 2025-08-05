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
	_ = godotenv.Load()

	db.Init()

	r := gin.Default()
	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from Go API!"})
	})

	// route public
	r.POST("/api/register", handler.Register)
	r.POST("/api/login",    handler.Login)

	r.GET("/debug-token", func(c *gin.Context) {
		tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		secret   := os.Getenv("JWT_SECRET")
	
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

	// route protected
	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired())
	
	{
		protected.POST("/invoices",        handler.CreateInvoice)
		protected.GET("/invoices",         handler.ListInvoices)
		protected.GET("/invoices/:id",     handler.GetInvoice)
		protected.PUT("/invoices/:id",     handler.UpdateInvoice)
		protected.DELETE("/invoices/:id",   handler.DeleteInvoice)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ""
	}
	log.Printf("Server ready → http://localhost:%s", port)
	r.Run(":" + port)
}