package db

import (
	"fmt"
	"log"
	"os"

	"invoice-api/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// getDatabaseDSN returns the database connection string from environment variables
func getDatabaseDSN() string {
	// Priority 1: Individual environment variables (lebih reliable untuk Render)
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	if host != "" && user != "" && password != "" && dbname != "" {
		if port == "" {
			port = "5432"
		}
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
			host, port, user, password, dbname)
		log.Printf("Using individual env vars connection to host: %s", host)
		return dsn
	}

	// Priority 2: DATABASE_URL (backup)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		log.Printf("Using DATABASE_URL connection")
		return dbURL
	}

	// Priority 3: POSTGRES_DSN (development)
	if dsn := os.Getenv("POSTGRES_DSN"); dsn != "" {
		log.Printf("Using POSTGRES_DSN connection")
		return dsn
	}

	// Fallback error
	log.Fatal("❌ No database configuration found. Please set individual DB environment variables")
	return ""
}

// maskPasswordInDSN masks password in connection string for logging
func maskPasswordInDSN(dsn string) string {
	if len(dsn) > 30 {
		return dsn[:30] + "***[MASKED]***"
	}
	return "***[MASKED]***"
}

func Init() {
	// Debug all database environment variables
	log.Printf("=== Database Environment Debug ===")
	log.Printf("DATABASE_URL: %s", maskPasswordInDSN(os.Getenv("DATABASE_URL")))
	log.Printf("DATABASE_HOST: %s", os.Getenv("DATABASE_HOST"))
	log.Printf("DATABASE_USER: %s", os.Getenv("DATABASE_USER"))
	log.Printf("DATABASE_NAME: %s", os.Getenv("DATABASE_NAME"))
	log.Printf("==================================")

	dsn := getDatabaseDSN()
	log.Printf("Final DSN: %s", maskPasswordInDSN(dsn))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ DB connect error: %v", err)
	}

	log.Println("✅ Database connection established successfully")

	// Test connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get underlying sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	// Auto migrate models
	log.Println("Running database migrations...")
	if err := DB.AutoMigrate(&model.User{}, &model.Invoice{}, &model.InvoiceItem{}); err != nil {
		log.Fatalf("❌ Failed to auto migrate: %v", err)
	}
	
	log.Println("✅ Database migrations completed successfully")
}