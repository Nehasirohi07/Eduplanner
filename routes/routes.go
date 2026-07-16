package routes

import (
	"eduplanner/handlers"
	"eduplanner/middleware"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
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

	router.Handle(
		"/students/{id}/goals",
		middleware.Auth(http.HandlerFunc(handlers.CreateStudyGoal)),
	).Methods("POST")

	router.Handle(
		"/subjects/{id}/goals",
		middleware.Auth(http.HandlerFunc(handlers.GetStudyGoals)),
	).Methods("GET")

	router.Handle(
		"/dashboard",
		middleware.Auth(http.HandlerFunc(handlers.GetDashboard)),
	).Methods("GET")

	router.PathPrefix("/swagger/").Handler(
		httpSwagger.WrapHandler,
	)

	return router
}
