package main

import (
	"github.com/nade-harlow/ecom-api/internal/adapter/database"
	"github.com/nade-harlow/ecom-api/internal/adapter/server"
)

func main() {
	db, _ := database.ConnectDb().DB()
	defer db.Close()

	database.RunManualMigration()

	server.StartUpServer()
}
