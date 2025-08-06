package db

import (
	"log"
	"os"
	"time"

	"invoice-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// maskPasswordInDSN masks password in connection string for logging
func maskPasswordInDSN(dsn string) string {
	if len(dsn) > 30 {
		return dsn[:30] + "***[MASKED]***"
	}
	return "***[MASKED]***"
}

// getDatabaseDSN returns the database URL from environment variable
func getDatabaseDSN() string {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL not set in environment")
	}
	log.Printf("Using DATABASE_URL: %s", maskPasswordInDSN(dbURL))
	return dbURL
}

// Init initializes the database connection and runs migrations
func Init() {
	log.Println("=== Database Initialization ===")

	dsn := getDatabaseDSN()
	log.Printf("Final DSN: %s", maskPasswordInDSN(dsn))

	var err error
	for attempts := 1; attempts <= 3; attempts++ {
		log.Printf("Database connection attempt %d/3", attempts)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Attempt %d failed: %v", attempts, err)
		if attempts < 3 {
			log.Println("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}

	if err != nil {
		log.Fatalf("❌ DB connect error after 3 attempts: %v", err)
	}

	log.Println("✅ Database connection established successfully")

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	log.Println("Running database migrations...")
	if err := DB.AutoMigrate(
		&model.User{},
		&model.Invoice{},
		&model.InvoiceItem{},
	); err != nil {
		log.Fatalf("❌ Failed to auto migrate: %v", err)
	}
	log.Println("✅ Database migrations completed successfully")
}
