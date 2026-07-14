package main

import (
	"fmt"
	"net/http"

	"eduplanner/config"
	"eduplanner/database"
	"eduplanner/routes"
)

func main() {
	config.LoadEnv()

	database.InitDB()

	database.CreateTables()

	router := routes.RegisterRoutes()

	fmt.Println("Server Starting...")

	err := http.ListenAndServe(":5050", router)

	if err != nil {
		fmt.Println(err)
	}
}
