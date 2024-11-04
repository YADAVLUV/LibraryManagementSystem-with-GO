package routes

import (
	"learningGorillamux/controllers"

	"github.com/gorilla/mux"
)

func UserRoutes(router *mux.Router) {
	router.HandleFunc("/login", controllers.Login)
	router.HandleFunc("/signup", controllers.Signup)
}
