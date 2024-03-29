package database

import (
	"chat-node/database/conversations"
	"chat-node/database/fetching"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB

func Connect() {
	url := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_DATABASE") + " port=" + os.Getenv("DB_PORT")

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
	db.AutoMigrate(&conversations.ConversationToken{})
	//db.AutoMigrate(&conversations.ConversationSpace{})
	db.AutoMigrate(&conversations.Message{})
	db.AutoMigrate(&fetching.Status{})

	// Assign the database to the global variable
	DBConn = db
}
