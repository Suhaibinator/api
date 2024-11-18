package main

import (
	"fmt"

	"os"

	"github.com/Suhaibinator/api/go_api"
	"github.com/Suhaibinator/api/go_persistence"
	"github.com/Suhaibinator/api/go_service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Database connection parameters
	username := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	hostname := os.Getenv("MYSQL_HOST")
	port := os.Getenv("port")
	dbname := os.Getenv("MYSQL_DATABASE")

	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, hostname, port, dbname)

	// Initialize the database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	appPersistence := go_persistence.NewApplicationPersistence(db)
	appService := go_service.NewApplicationService(appPersistence)
	appRouter := go_api.NewApplicationRouter(appService)

	appRouter.RegisterRoutes()
	appRouter.Run(8084)

}
