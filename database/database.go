package database

import (
	"chat-node/database/conversations"
	"chat-node/database/fetching"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB

func Connect() {
	url := "host=localhost user=postgres password=deinemutter123 dbname=chat port=5432"

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatal("Something went wrong during the connection with the database.", err)
	}

	log.Println("Successfully connected to the database.")

	// Configure the database driver
	driver, _ := db.DB()

	driver.SetMaxIdleConns(10)
	driver.SetMaxOpenConns(100)
	driver.SetConnMaxLifetime(time.Hour)

	// Migrate the schema
	db.AutoMigrate(&conversations.Conversation{})
	db.AutoMigrate(&conversations.Member{})
	db.AutoMigrate(&conversations.Message{})
	db.AutoMigrate(&fetching.Session{})

	// Assign the database to the global variable
	DBConn = db
}
