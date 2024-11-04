package routes

import (
	"learningGorillamux/controllers"
	"learningGorillamux/middleware"

	"github.com/gorilla/mux"
)

func BookRoutes(router *mux.Router) {
	router.HandleFunc("/getBooks", controllers.GetBooks)
	router.HandleFunc("/getBook/{bookId}", controllers.GetBook)
	router.HandleFunc("/deletebook/{bookId}", middleware.ValidateOwner(controllers.DeleteBook))
	router.HandleFunc("/addBook", middleware.ValidateOwner(controllers.AddBook))
}
