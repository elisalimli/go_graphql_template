package main

import (
	"log"
	"os"

	"github.com/elisalimli/go_graphql_template/graphql/models"
	"github.com/elisalimli/go_graphql_template/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main() {
	log.Println("Running Migrations")
	err := initializers.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Migration Failed:  \n", err.Error())
		os.Exit(1)
	}

}
