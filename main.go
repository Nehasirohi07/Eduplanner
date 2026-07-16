package main

import (
	"fmt"
	"net/http"

	"eduplanner/config"
	"eduplanner/database"
	_ "eduplanner/docs"
	"eduplanner/routes"
)

//@title eduPlanner API
//@version 1.0
//@description Backend API for EduPlanner
//@host localhost:5050
//@BasePath /

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
