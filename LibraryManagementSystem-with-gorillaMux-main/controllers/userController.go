package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"learningGorillamux/database"
	"learningGorillamux/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Login(w http.ResponseWriter, r *http.Request) {

	var user models.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &user)

	var user_in_db models.User
	err = database.UserCollection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&user_in_db)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user was not found"))
		return
	}
	if user_in_db.Password == user.Password {
		claims := models.Claims{
			Username: user.Username,
			UserType: user.UserType,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		encTkn, err := token.SignedString(models.JwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusFound)
		w.Write([]byte(encTkn))
		return
	}
	w.Write([]byte("wrong password"))
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var newUser models.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &newUser)

	if err != nil {
		http.Error(w, "Error while unmarsheling data", http.StatusInternalServerError)
		return
	}

	existingUser := database.UserCollection.FindOne(context.Background(), bson.M{"username": newUser.Username})
	if existingUser.Err() == nil {
		w.Write([]byte("User already exists!"))
		return
	}

	newUserToInsertInDb := models.User{
		Id:         primitive.NewObjectID(),
		Username:   newUser.Username,
		Password:   newUser.Password,
		UserType:   newUser.UserType,
		Created_at: time.Now(),
	}
	_, err = database.UserCollection.InsertOne(context.Background(), newUserToInsertInDb)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("new user created with userId %s and password %s", newUser.Username, newUser.Password)))

}