package controllers

import (
	"context"
	"encoding/json"
	"io"
	"learningGorillamux/database"
	"learningGorillamux/models"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookIdstr := vars["bookId"]

	var RequiredBook models.Book

	id, _ := primitive.ObjectIDFromHex(bookIdstr)
	filter := bson.M{"_id": id}
	err := database.BookCollection.FindOne(context.Background(), filter).Decode(&RequiredBook)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Book Doesn't exist"))
		return
	}

	marsheledBook, err := json.Marshal(RequiredBook)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	w.Write(marsheledBook)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {

	coursor, err := database.BookCollection.Find(context.Background(), bson.D{})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No Books Available"))
		return
	}

	var listOfBooks models.Books

	for coursor.Next(context.Background()) {
		var book models.Book
		err := coursor.Decode(&book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		listOfBooks.AddBookToList(book)
	}

	encData, err := json.Marshal(listOfBooks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while marsheling the data"))
		return
	}

	w.WriteHeader(http.StatusFound)
	w.Header().Set("Content-Type", "application/json")
	w.Write(encData)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookIdstr := vars["bookId"]

	id, _ := primitive.ObjectIDFromHex(bookIdstr)
	filter := bson.M{"_id": id}
	_, err := database.BookCollection.DeleteOne(context.Background(), filter)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Book Doesn't exist"))
		return
	}

	w.Write([]byte("Successfully deleted Book!"))
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &newBook)

	if err != nil {
		http.Error(w, "Error while unmarsheling data", http.StatusInternalServerError)
		return
	}

	existingBook := database.BookCollection.FindOne(context.Background(), bson.M{"book_name": newBook.Book_name})
	if existingBook.Err() == nil {
		w.Write([]byte("User already exists!"))
		return
	}

	newBookToAddInDb := models.Book{
		Id:         primitive.NewObjectID(),
		Book_name:  newBook.Book_name,
		Author:     newBook.Author,
		Created_at: time.Now(),
	}

	_, err = database.BookCollection.InsertOne(context.Background(), newBookToAddInDb)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("New Book Added to DB"))
}
