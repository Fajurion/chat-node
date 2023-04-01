package database

import (
	"chat-node/database/conversations"
	"chat-node/database/credentials"
	"chat-node/database/fetching"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB

func Connect() {
	url := "host=" + credentials.DB_HOST + " user=" + credentials.DB_USERNAME + " password=" + credentials.DB_PASSWORD + " dbname=" + credentials.DB_DATABASE + " port=" + credentials.DB_PORT

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
	db.AutoMigrate(&fetching.Status{})
	db.AutoMigrate(&fetching.Action{})

	// Assign the database to the global variable
	DBConn = db
}
