package main

import (
	"learningGorillamux/middleware"
	"learningGorillamux/routes"
	"net/http"

	"github.com/gorilla/mux"
)



func main() {

	router := mux.NewRouter()

	routes.BookRoutes(router)
	routes.UserRoutes(router)
	routes.OrderRoutes(router)

	router.Use(middleware.TrackNumberOfRequests)

	http.ListenAndServe(":8000", router)
}
