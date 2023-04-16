package initializers

import (
	"log"
	"os"

	"github.com/elisalimli/go_graphql_template/graphql/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectToDatabase() {
	var err error
	dsn := os.Getenv("POSTGRESQL_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	DB.Logger = logger.Default.LogMode(logger.Info)

	DB.AutoMigrate(&models.User{})

	log.Println("🚀 Connected Successfully to the Database")

	if err != nil {
		log.Fatal("Failed to connect a database.")
	}
}
