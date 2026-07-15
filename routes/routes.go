package routes

import (
	"eduplanner/handlers"
	"eduplanner/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/register", handlers.Register).Methods("POST")

	router.HandleFunc("/login", handlers.Login).Methods("POST")

	router.Handle(
		"/courses",
		middleware.Auth(http.HandlerFunc(handlers.CreateCourse)),
	).Methods("POST")

	router.Handle(
		"/courses/{id}/subjects",
		middleware.Auth(http.HandlerFunc(handlers.CreateSubject)),
	).Methods("POST")

	router.Handle(
		"/courses/{id}/subjects",
		middleware.Auth(http.HandlerFunc(handlers.GetSubjects)),
	).Methods("GET")
	return router

}
