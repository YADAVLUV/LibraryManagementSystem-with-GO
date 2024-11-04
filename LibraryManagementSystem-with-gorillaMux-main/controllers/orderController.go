package controllers

import (
	"context"
	"encoding/json"
	"learningGorillamux/database"
	"learningGorillamux/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func OrderBook(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("JwtToken")
	if tokenStr == "" {
		w.Write([]byte("User not logged In"))
		return
	}
	claims := &models.Claims{}

	jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return models.JwtKey, nil
	})

	username := claims.Username
	vars := mux.Vars(r)
	bookIdstr, errbool := vars["bookId"]
	if !errbool {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Book id not provided in url"))
		return
	}
	id, _ := primitive.ObjectIDFromHex(bookIdstr)
	filter := bson.M{"_id": id}
	BookExists := database.BookCollection.FindOne(context.Background(), filter)
	if BookExists.Err() != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Book with this id doesn't exist"))
		return
	}

	order_book := models.Order{
		Id:         primitive.NewObjectID(),
		User_id:    username,
		Book_id:    bookIdstr,
		Created_at: time.Now(),
	}

	_, err := database.OrderCollection.InsertOne(context.Background(), order_book)
	if err != nil {
		panic(err)
	}

	w.Write([]byte("Successfully ordered book!"))
}

func ListAllOrderedBooks(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("JwtToken")
	if tokenStr == "" {
		w.Write([]byte("User not logged In"))
		return
	}
	claims := &models.Claims{}

	jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return models.JwtKey, nil
	})
	username := claims.Username
	////
	coursor, err := database.OrderCollection.Find(context.Background(), bson.M{"user_id": username})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No Books Available"))
		return
	}

	var listOfOrders models.Orders

	for coursor.Next(context.Background()) {
		var order models.Order
		err := coursor.Decode(&order)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		listOfOrders.AddOrderToList(order)
	}

	encData, err := json.Marshal(listOfOrders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while marsheling the data"))
		return
	}

	w.WriteHeader(http.StatusFound)
	w.Header().Set("Content-Type", "application/json")
	w.Write(encData)
}
